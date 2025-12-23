package email

import (
	"fmt"
	"io"
	"lain/types"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

func SelectFolder(client *types.EmailClient, folderName string) (*imap.MailboxStatus, error) {
	mbox, err := client.Select(folderName, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select folder: %w", err)
	}
	return mbox, nil
}

func FetchMessages(client *types.EmailClient, folderName string, limit uint32) ([]*types.EmailMessage, error) {
	mbox, err := SelectFolder(client, folderName)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return []*types.EmailMessage{}, nil
	}

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > limit {
		from = mbox.Messages - limit + 1
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid, imap.FetchRFC822Size, section.FetchItem()}

	go func() {
		done <- client.Fetch(seqset, items, messages)
	}()

	var result []*types.EmailMessage
	for msg := range messages {
		if msg == nil {
			continue
		}

		emailMsg, err := parseMessage(msg)
		if err != nil {
			continue
		}

		result = append(result, emailMsg)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	return result, nil
}

func parseMessage(msg *imap.Message) (*types.EmailMessage, error) {
	if msg.Envelope == nil {
		return nil, fmt.Errorf("message envelope is nil")
	}

	section := &imap.BodySectionName{}
	bodyReader := msg.GetBody(section)
	if bodyReader == nil {
		return nil, fmt.Errorf("message body is nil")
	}

	mr, err := mail.CreateReader(bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail reader: %w", err)
	}

	var bodyText, bodyHTML string
	var attachments []types.EmailAttachment

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			body, _ := io.ReadAll(part.Body)

			if contentType == "text/plain" {
				bodyText = string(body)
			} else if contentType == "text/html" {
				bodyHTML = string(body)
			}

		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			contentType, _, _ := h.ContentType()
			data, _ := io.ReadAll(part.Body)

			attachments = append(attachments, types.EmailAttachment{
				Filename:    filename,
				ContentType: contentType,
				Data:        data,
			})
		}
	}

	var fromAddr, fromName string
	if len(msg.Envelope.From) > 0 {
		fromAddr = msg.Envelope.From[0].MailboxName + "@" + msg.Envelope.From[0].HostName
		fromName = msg.Envelope.From[0].PersonalName
	}

	var toList []string
	for _, addr := range msg.Envelope.To {
		toList = append(toList, addr.MailboxName+"@"+addr.HostName)
	}

	var ccList []string
	for _, addr := range msg.Envelope.Cc {
		ccList = append(ccList, addr.MailboxName+"@"+addr.HostName)
	}

	var bccList []string
	for _, addr := range msg.Envelope.Bcc {
		bccList = append(bccList, addr.MailboxName+"@"+addr.HostName)
	}

	var replyToList []string
	for _, addr := range msg.Envelope.ReplyTo {
		replyToList = append(replyToList, addr.MailboxName+"@"+addr.HostName)
	}

	isRead := false
	isFlagged := false
	isAnswered := false
	isDraft := false

	for _, flag := range msg.Flags {
		switch flag {
		case imap.SeenFlag:
			isRead = true
		case imap.FlaggedFlag:
			isFlagged = true
		case imap.AnsweredFlag:
			isAnswered = true
		case imap.DraftFlag:
			isDraft = true
		}
	}

	return &types.EmailMessage{
		UID:           msg.Uid,
		MessageID:     msg.Envelope.MessageId,
		From:          fromAddr,
		FromName:      fromName,
		To:            toList,
		CC:            ccList,
		BCC:           bccList,
		ReplyTo:       replyToList,
		Subject:       msg.Envelope.Subject,
		Date:          msg.Envelope.Date,
		BodyText:      bodyText,
		BodyHTML:      bodyHTML,
		Size:          msg.Size,
		InReplyTo:     msg.Envelope.InReplyTo,
		IsRead:        isRead,
		IsFlagged:     isFlagged,
		IsAnswered:    isAnswered,
		IsDraft:       isDraft,
		HasAttachment: len(attachments) > 0,
		Attachments:   attachments,
	}, nil
}

func MarkAsRead(client *types.EmailClient, folderName string, uid uint32) error {
	if _, err := SelectFolder(client, folderName); err != nil {
		return err
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uid)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}

	if err := client.UidStore(seqSet, item, flags, nil); err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	return nil
}

func ToggleFlag(client *types.EmailClient, folderName string, uid uint32, isFlagged bool) error {
	if _, err := SelectFolder(client, folderName); err != nil {
		return err
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uid)

	var item imap.StoreItem
	if isFlagged {
		item = imap.FormatFlagsOp(imap.RemoveFlags, true)
	} else {
		item = imap.FormatFlagsOp(imap.AddFlags, true)
	}

	flags := []interface{}{imap.FlaggedFlag}

	if err := client.UidStore(seqSet, item, flags, nil); err != nil {
		return fmt.Errorf("failed to toggle flag: %w", err)
	}

	return nil
}
