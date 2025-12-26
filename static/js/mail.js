document.addEventListener('DOMContentLoaded', function () {
    const emailRows = document.querySelectorAll('.email-row');
    const preview = document.querySelector('.preview');
    const prefsData = document.getElementById('mail-preferences');

    let currentEmailId = null;
    let markAsReadTimer = null;

    // Parse preferences from data attributes
    const prefs = {
        MarkMessagesAsRead: prefsData ? prefsData.dataset.markAsRead : 'Immediately',
        ShowEmailAddressWithDisplayName: prefsData ? prefsData.dataset.showAddress === 'true' : true,
        DisplayHTML: prefsData ? prefsData.dataset.displayHtml === 'true' : true,
        LoadRemoteContent: prefsData ? prefsData.dataset.loadRemote : 'Never'
    };

    emailRows.forEach(row => {
        row.addEventListener('click', async function (e) {
            if (e.target.closest('.email-flag')) {
                return;
            }

            const emailId = this.dataset.emailId;
            if (emailId === currentEmailId) {
                return;
            }

            emailRows.forEach(r => r.classList.remove('active'));
            this.classList.add('active');

            currentEmailId = emailId;

            // Clear previous timer
            if (markAsReadTimer) {
                clearTimeout(markAsReadTimer);
            }

            try {
                const response = await fetch(`/api/mail/email/${emailId}`);
                if (!response.ok) throw new Error('Failed to fetch email');

                const email = await response.json();
                renderEmail(email);

                // Handle mark as read based on preference
                if (!email.IsRead) {
                    handleMarkAsRead(emailId, this);
                }
            } catch (error) {
                console.error('Error fetching email:', error);
                showError('Error loading email');
            }
        });
    });

    // Handle flag clicks
    document.querySelectorAll('.email-flag').forEach(flag => {
        flag.addEventListener('click', async function (e) {
            e.stopPropagation();

            const row = this.closest('.email-row');
            const emailId = row.dataset.emailId;

            try {
                const response = await fetch(`/api/mail/email/${emailId}/flag`, {
                    method: 'POST'
                });

                if (!response.ok) throw new Error('Failed to toggle flag');

                const data = await response.json();

                if (data.flagged) {
                    this.classList.add('flagged');
                    this.title = 'Unflag';
                } else {
                    this.classList.remove('flagged');
                    this.title = 'Flag';
                }
            } catch (error) {
                console.error('Error toggling flag:', error);
            }
        });
    });

    function handleMarkAsRead(emailId, row) {
        const markOption = prefs.MarkMessagesAsRead || 'Immediately';

        const delays = {
            'Never': null,
            'Immediately': 0,
            'After 5 Seconds': 5000,
            'After 10 Seconds': 10000,
            'After 30 Seconds': 30000,
            'After 1 Minute': 60000
        };

        const delay = delays[markOption];

        if (delay === null) {
            return; // Never mark as read
        }

        markAsReadTimer = setTimeout(async () => {
            try {
                const response = await fetch(`/api/mail/email/${emailId}/read`, {
                    method: 'POST'
                });

                if (response.ok) {
                    row.classList.remove('unread');
                }
            } catch (error) {
                console.error('Error marking as read:', error);
            }
        }, delay);
    }

    function renderEmail(email) {
        preview.innerHTML = '';

        // Header
        const header = createHeader(email);
        preview.appendChild(header);

        // Sender info
        const sender = createSenderInfo(email);
        preview.appendChild(sender);

        // Recipients
        const recipients = createRecipients(email);
        preview.appendChild(recipients);

        // Attachments
        if (email.Attachments && email.Attachments.length > 0) {
            const attachments = createAttachments(email.Attachments);
            preview.appendChild(attachments);
        }

        // Body
        const body = createBody(email);
        preview.appendChild(body);
    }

    function createHeader(email) {
        const header = document.createElement('div');
        header.className = 'email-header';

        const subject = document.createElement('h2');
        subject.className = 'email-subject';

        if (email.Subject) {
            subject.textContent = email.Subject;
        } else {
            const noSubject = document.createElement('span');
            noSubject.className = 'no-subject';
            noSubject.textContent = '[No Subject]';
            subject.appendChild(noSubject);
        }

        const actions = document.createElement('div');
        actions.className = 'email-actions';

        const actionButtons = [
            { title: 'Reply', symbol: '↶' },
            { title: 'Reply All', symbol: '⇄' },
            { title: 'Forward', symbol: '→' },
            { title: 'Archive', symbol: '▼' },
            { title: 'Delete', symbol: '×' }
        ];

        actionButtons.forEach(btn => {
            const button = document.createElement('button');
            button.className = 'btn-icon';
            button.title = btn.title;
            button.textContent = btn.symbol;
            actions.appendChild(button);
        });

        header.appendChild(subject);
        header.appendChild(actions);

        return header;
    }

    function createSenderInfo(email) {
        const sender = document.createElement('div');
        sender.className = 'email-sender';

        const senderInfo = document.createElement('div');
        senderInfo.className = 'sender-info';

        // Respect ShowEmailAddressWithDisplayName preference
        const showAddress = prefs.ShowEmailAddressWithDisplayName;

        const strong = document.createElement('strong');
        strong.textContent = email.FromName || email.From;
        senderInfo.appendChild(strong);

        if (showAddress && email.FromName) {
            const address = document.createTextNode(` <${email.From}>`);
            senderInfo.appendChild(address);
        }

        const dateDiv = document.createElement('div');
        dateDiv.className = 'email-date';
        dateDiv.textContent = formatDate(email.Date);

        sender.appendChild(senderInfo);
        sender.appendChild(dateDiv);

        return sender;
    }

    function createRecipients(email) {
        const recipients = document.createElement('div');
        recipients.className = 'email-recipients';

        const toDiv = document.createElement('div');
        const toStrong = document.createElement('strong');
        toStrong.textContent = 'To: ';
        toDiv.appendChild(toStrong);
        toDiv.appendChild(document.createTextNode(email.To || ''));
        recipients.appendChild(toDiv);

        if (email.CC) {
            const ccDiv = document.createElement('div');
            const ccStrong = document.createElement('strong');
            ccStrong.textContent = 'Cc: ';
            ccDiv.appendChild(ccStrong);
            ccDiv.appendChild(document.createTextNode(email.CC));
            recipients.appendChild(ccDiv);
        }

        return recipients;
    }

    function createAttachments(attachments) {
        const container = document.createElement('div');
        container.className = 'email-attachments';

        const label = document.createElement('strong');
        label.textContent = 'Attachments: ';
        container.appendChild(label);

        attachments.forEach(att => {
            const link = document.createElement('a');
            link.href = `/api/mail/attachment/${att.ID}`;
            link.className = 'attachment';
            link.download = att.Filename;
            link.textContent = `${att.Filename} (${att.Size})`;
            container.appendChild(link);
        });

        return container;
    }

    function createBody(email) {
        const body = document.createElement('div');
        body.className = 'email-body';
        const displayHTML = prefs.DisplayHTML;

        if (displayHTML && email.Body) {
            const shadow = ShadowRenderer.render(body, email.Body);
            handleRemoteContent(shadow);
        } else {
            const pre = document.createElement('pre');
            pre.textContent = email.Body || '[Empty Content]';
            body.appendChild(pre);
        }

        return body;
    }

    function handleRemoteContent(container) {
        const loadOption = prefs.LoadRemoteContent || 'Never';

        if (loadOption === 'Never') {
            // Block all external images
            const images = container.querySelectorAll('img');
            images.forEach(img => {
                const src = img.getAttribute('src');
                if (src && src.startsWith('http')) {
                    img.removeAttribute('src');
                    img.dataset.src = src; // Store original src
                    img.alt = '[Remote image blocked]';
                    img.style.border = '1px dashed #ccc';
                    img.style.padding = '5px';
                    img.style.display = 'inline-block';
                }
            });
        }
        // TODO: Implement "From my contacts" check
        // For "Always", images load normally
    }

    function showError(message) {
        preview.innerHTML = '';
        const error = document.createElement('div');
        error.className = 'no-email-selected';
        const p = document.createElement('p');
        p.textContent = message;
        error.appendChild(p);
        preview.appendChild(error);
    }

    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString('en-US', {
            weekday: 'short',
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }
});