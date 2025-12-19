package types

type EmailReadingOption string

const (
	EmailReadingOptionNever          EmailReadingOption = "Never"
	EmailReadingOptionImmediately    EmailReadingOption = "Immediately"
	EmailReadingOptionAfter5Seconds  EmailReadingOption = "After 5 Seconds"
	EmailReadingOptionAfter10Seconds EmailReadingOption = "After 10 Seconds"
	EmailReadingOptionAfter30Seconds EmailReadingOption = "After 30 Seconds"
	EmailReadingOptionAfter1Minute   EmailReadingOption = "After 1 Minute"
)

type EmailComposingOption string

const (
	EmailComposingOptionNeverHTML                           EmailComposingOption = "Never"
	EmailComposingOptionOnReplyToHTML                       EmailComposingOption = "When replying to HTML messages"
	EmailComposingOptionOnForwardOrReplyToHTML              EmailComposingOption = "When forwarding or replying to HTML messages"
	EmailComposingOptionAlwaysExceptWhenReplyingToPlainText EmailComposingOption = "Always, except when replying to plain text messages"
	EmailComposingOptionAlwaysHTML                          EmailComposingOption = "Always"
)

type RemoteResourceDownloadOption string

const (
	RemoteResourceDownloadOptionNever          RemoteResourceDownloadOption = "Never"
	RemoteResourceDownloadOptionFromMyContacts RemoteResourceDownloadOption = "From my contacts"
	RemoteResourceDownloadOptionAlways         RemoteResourceDownloadOption = "Always"
)

type ReturnReceiptOption string

const (
	ReturnReceiptOptionAskMe                         ReturnReceiptOption = "Ask me each time"
	ReturnReceiptOptionSendAlways                    ReturnReceiptOption = "Always send a receipt"
	ReturnReceiptOptionIgnore                        ReturnReceiptOption = "Ignore all requests"
	ReturnReceiptOptionSendToContactsOtherwiseAsk    ReturnReceiptOption = "Send receipt to my contacts, otherwise ask me"
	ReturnReceiptOptionSendToContactsOtherwiseIgnore ReturnReceiptOption = "Send receipt to my contacts, otherwise ignore"
)

type AutoSaveDraftIntervalOption int

const (
	AutoSaveDraftIntervalOptionNever     AutoSaveDraftIntervalOption = 0
	AutoSaveDraftIntervalOption30Seconds AutoSaveDraftIntervalOption = 30
	AutoSaveDraftIntervalOption1Minute   AutoSaveDraftIntervalOption = 60
	AutoSaveDraftIntervalOption3Minutes  AutoSaveDraftIntervalOption = 180
	AutoSaveDraftIntervalOption5Minutes  AutoSaveDraftIntervalOption = 300
	AutoSaveDraftIntervalOption10Minutes AutoSaveDraftIntervalOption = 600
)

type EmailReplyOption string

const (
	EmailReplyOptionDoNotQuote EmailReplyOption = "Do not quote the original message"
	EmailReplyOptionBelowQuote EmailReplyOption = "Place my reply below the original message"
	EmailReplyOptionAboveQuote EmailReplyOption = "Place my reply above the original message"
)

type MessageForwardingOption string

const (
	MessageForwardingOptionAsAttachment MessageForwardingOption = "As attachment"
	MessageForwardingOptionInline       MessageForwardingOption = "Inline"
)

type EmailHTMLFontFamilyOption string

const (
	EmailHTMLFontFamilyOptionAndaleMono    EmailHTMLFontFamilyOption = "Andale Mono"
	EmailHTMLFontFamilyOptionArial         EmailHTMLFontFamilyOption = "Arial"
	EmailHTMLFontFamilyOptionArialBlack    EmailHTMLFontFamilyOption = "Arial Black"
	EmailHTMLFontFamilyOptionBookAntiqua   EmailHTMLFontFamilyOption = "Book Antiqua"
	EmailHTMLFontFamilyOptionComicSansMS   EmailHTMLFontFamilyOption = "Comic Sans MS"
	EmailHTMLFontFamilyOptionCourierNew    EmailHTMLFontFamilyOption = "Courier New"
	EmailHTMLFontFamilyOptionGeorgia       EmailHTMLFontFamilyOption = "Georgia"
	EmailHTMLFontFamilyOptionHelvetica     EmailHTMLFontFamilyOption = "Helvetica"
	EmailHTMLFontFamilyOptionImpact        EmailHTMLFontFamilyOption = "Impact"
	EmailHTMLFontFamilyOptionTahoma        EmailHTMLFontFamilyOption = "Tahoma"
	EmailHTMLFontFamilyOptionTerminal      EmailHTMLFontFamilyOption = "Terminal"
	EmailHTMLFontFamilyOptionTimesNewRoman EmailHTMLFontFamilyOption = "Times New Roman"
	EmailHTMLFontFamilyOptionTrebuchetMS   EmailHTMLFontFamilyOption = "Trebuchet MS"
	EmailHTMLFontFamilyOptionVerdana       EmailHTMLFontFamilyOption = "Verdana"
)

type EmailHTMLFontSizeOption int

const (
	EmailHTMLFontSizeOption8Pt  EmailHTMLFontSizeOption = 8
	EmailHTMLFontSizeOption9Pt  EmailHTMLFontSizeOption = 9
	EmailHTMLFontSizeOption10Pt EmailHTMLFontSizeOption = 10
	EmailHTMLFontSizeOption11Pt EmailHTMLFontSizeOption = 11
	EmailHTMLFontSizeOption12Pt EmailHTMLFontSizeOption = 12
	EmailHTMLFontSizeOption14Pt EmailHTMLFontSizeOption = 14
	EmailHTMLFontSizeOption16Pt EmailHTMLFontSizeOption = 16
	EmailHTMLFontSizeOption18Pt EmailHTMLFontSizeOption = 18
	EmailHTMLFontSizeOption24Pt EmailHTMLFontSizeOption = 24
	EmailHTMLFontSizeOption36Pt EmailHTMLFontSizeOption = 36
)

type EmailSignatureOption string

const (
	EmailSignatureOptionAlways         EmailSignatureOption = "Always"
	EmailSignatureOptionNever          EmailSignatureOption = "Never"
	EmailSignatureOptionForNewMessages EmailSignatureOption = "For new messages only"
	EmailSignatureOptionForReplies     EmailSignatureOption = "For replies and forwards only"
)

type EmailAttachementNameStyleOption string

const (
	EmailAttachementNameStyleOptionThunderbird EmailAttachementNameStyleOption = "Full RFC 2231 (Thunderbird)"
	EmailAttachementNameStyleOptionOutlook     EmailAttachementNameStyleOption = "RFC 2047/2231 (Outlook)"
	EmailAttachementNameStyleOptionOther       EmailAttachementNameStyleOption = "Full RFC 2231 (Other)"
)

type ContactDisplayOption string

const (
	ContactDisplayOptionDisplayName             ContactDisplayOption = "Display Name"
	ContactDisplayOptionFirstLast               ContactDisplayOption = "First Last"
	ContactDisplayOptionLastFirst               ContactDisplayOption = "Last First"
	ContactDisplayOptionLastFirstCommaSeparated ContactDisplayOption = "Last, First"
)

type ContactSortingOption string

const (
	ContactSortingOptionFirstName   ContactSortingOption = "First Name"
	ContactSortingOptionLastName    ContactSortingOption = "Last Name"
	ContactSortingOptionDisplayName ContactSortingOption = "Display Name"
)

type ClearTrashOnLogoutOption string

const (
	ClearTrashOnLogoutOptionNever           ClearTrashOnLogoutOption = "Never"
	ClearTrashOnLogoutOptionAllMessages     ClearTrashOnLogoutOption = "All messages"
	ClearTrashOnLogoutOptionOlderThan30Days ClearTrashOnLogoutOption = "Messages older than 30 days"
	ClearTrashOnLogoutOptionOlderThan60Days ClearTrashOnLogoutOption = "Messages older than 60 days"
	ClearTrashOnLogoutOptionOlderThan90Days ClearTrashOnLogoutOption = "Messages older than 90 days"
)
