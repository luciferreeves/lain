document.addEventListener('DOMContentLoaded', function () {
    const emailRows = document.querySelectorAll('.email-row');
    const preview = document.querySelector('.preview');
    const prefsData = document.getElementById('mail-preferences');

    let currentEmailId = null;
    let markAsReadTimer = null;
    let currentEmail = null;
    let viewMode = 'html';
    let wordWrap = true;
    let showDetails = false;

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

            if (markAsReadTimer) {
                clearTimeout(markAsReadTimer);
            }

            try {
                const response = await fetch(`/api/mail/email/${emailId}`);
                if (!response.ok) throw new Error('Failed to fetch email');

                const email = await response.json();
                currentEmail = email;
                showDetails = false;
                renderEmail(email);

                document.querySelectorAll('.subnav .nav-subitem').forEach(item => {
                    item.removeAttribute('disabled');
                    item.classList.remove('disabled');
                });

                if (!email.IsRead) {
                    handleMarkAsRead(emailId, this);
                }
            } catch (error) {
                console.error('Error fetching email:', error);
                showError('Error loading email');
            }
        });
    });

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

        if (delay === null) return;

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

    async function renderEmail(email) {
        preview.innerHTML = '';

        const header = await createHeader(email);
        preview.appendChild(header);

        const body = createBody(email);
        preview.appendChild(body);
    }

    async function createHeader(email) {
        const header = document.createElement('div');
        header.className = 'email-card-header';

        // Card container
        const card = document.createElement('div');
        card.className = 'email-card';

        // Header Top Row (Subject + Actions)
        const headerRow = document.createElement('div');
        headerRow.className = 'header-top-row';

        // Subject container
        const subjectContainer = document.createElement('div');
        subjectContainer.className = 'subject-container';
        subjectContainer.innerHTML = `
            <div class="field-label">SUBJECT</div>
            <div class="field-value subject-text">${email.Subject || '[No Subject]'}</div>
        `;

        // Action buttons grid
        const actionsGrid = document.createElement('div');
        actionsGrid.className = 'card-actions-grid';

        const actions = [
            { icon: Icons.Html, label: 'Switch to HTML View', active: viewMode === 'html', onClick: () => switchViewMode('html') },
            { icon: Icons.Plain, label: 'Switch to Plain Text View', active: viewMode === 'plain', onClick: () => switchViewMode('plain') },
            { icon: Icons.Reply, label: 'Reply to Sender', onClick: () => console.log('Reply') },
            { icon: Icons.Forward, label: 'Forward Message', onClick: () => console.log('Forward') },
            {
                id: 'btn-details',
                icon: showDetails ? Icons.Summary : Icons.Details,
                label: showDetails ? 'Hide Details' : 'Show Details',
                onClick: toggleDetails
            },
            { icon: Icons.Wrap, label: 'Toggle Word Wrap', active: wordWrap, onClick: () => toggleWordWrap() },
            { icon: Icons.Headers, label: 'View Message Headers', onClick: () => showHeaders(email) },
            { icon: Icons.Window, label: 'Open in New Window', onClick: () => console.log('Window') }
        ];

        actions.forEach(action => {
            const btn = document.createElement('button');
            btn.className = 'card-action-btn';
            if (action.id) btn.id = action.id;
            if (action.active) btn.classList.add('active');
            btn.innerHTML = action.icon;
            btn.title = action.label; // Tooltip
            btn.onclick = action.onClick;
            actionsGrid.appendChild(btn);
        });

        headerRow.appendChild(subjectContainer);
        headerRow.appendChild(actionsGrid);

        // From field with profile
        const fromField = document.createElement('div');
        fromField.className = 'card-field card-field-with-pic';

        const picUrl = await EmailUtils.getProfilePicture(email.FromEmail, email.FromName);
        const escapedName = (email.FromName || '').replace(/'/g, "\\'");
        const escapedEmail = (email.FromEmail || '').replace(/'/g, "\\'");

        fromField.innerHTML = `
        <div class="field-label">FROM</div>
        <div class="field-value-with-pic">
            <img src="${picUrl}" class="card-pic" onerror="this.outerHTML='<div class=card-pic-init>${EmailUtils.getInitials(email.FromName, email.FromEmail)}</div>'">
            <div class="sender-details">
                <div class="sender-name-card" onclick="showSenderMenu(event, '${escapedName}', '${escapedEmail}')">${email.FromName || email.FromEmail}</div>
                <div class="sender-email-card" onclick="showSenderMenu(event, '${escapedName}', '${escapedEmail}')">${email.FromEmail}</div>
            </div>
        </div>
        `;

        // Date field
        const dateField = document.createElement('div');
        dateField.className = 'card-field';
        dateField.innerHTML = `
            <div class="field-label">DATE</div>
            <div class="field-value">${email.DateFormatted}</div>
        `;

        // Details container (hidden by default)
        const detailsContainer = document.createElement('div');
        detailsContainer.id = 'email-details';
        detailsContainer.style.display = 'none';

        // TO field
        const toRow = document.createElement('div');
        toRow.className = 'card-field';
        toRow.innerHTML = `
            <div class="field-label">TO</div>
            <div class="field-value">${email.To}</div>
        `;
        detailsContainer.appendChild(toRow);

        // CC field
        if (email.CC) {
            const ccRow = document.createElement('div');
            ccRow.className = 'card-field';
            ccRow.innerHTML = `
                <div class="field-label">CC</div>
                <div class="field-value">${email.CC}</div>
            `;
            detailsContainer.appendChild(ccRow);
        }

        card.appendChild(headerRow);
        card.appendChild(fromField);
        card.appendChild(dateField);
        card.appendChild(detailsContainer);

        // Attachments field
        if (email.Attachments && email.Attachments.length > 0) {
            const attRow = document.createElement('div');
            attRow.className = 'card-field';

            let attHtml = '';
            email.Attachments.forEach(att => {
                attHtml += `<a href="/api/mail/attachment/${att.ID}" class="attachment" download="${att.Filename}">${att.Filename} (${att.Size})</a>`;
            });

            attRow.innerHTML = `
                <div class="field-label">ATTACHMENTS</div>
                <div class="field-value">${attHtml}</div>
            `;
            card.appendChild(attRow);
        }

        header.appendChild(card);
        return header;
    }



    function toggleDetails() {
        showDetails = !showDetails;

        const detailsContainer = document.getElementById('email-details');
        const detailsBtn = document.getElementById('btn-details');

        if (detailsContainer) {
            detailsContainer.style.display = showDetails ? 'block' : 'none';
        }

        if (detailsBtn) {
            const icon = showDetails ? Icons.Summary : Icons.Details;
            const label = showDetails ? 'Hide Details' : 'Show Details';
            detailsBtn.innerHTML = icon;
            detailsBtn.title = label;
        }
    }

    function showSenderMenu(e, name, email) {
        const existingMenu = document.querySelector('.sender-menu');
        if (existingMenu) existingMenu.remove();

        const menu = document.createElement('div');
        menu.className = 'sender-menu';

        const addContact = document.createElement('div');
        addContact.className = 'sender-menu-item';
        addContact.innerHTML = Icons.addContact + ' <span>Add to Address Book</span>';
        addContact.onclick = () => {
            console.log('Add to address book:', name, email);
            menu.remove();
        };

        const composeMail = document.createElement('div');
        composeMail.className = 'sender-menu-item';
        composeMail.innerHTML = Icons.composeMail + ' <span>Compose Mail to</span>';
        composeMail.onclick = () => {
            console.log('Compose mail to:', name, email);
            menu.remove();
        };

        menu.appendChild(addContact);
        menu.appendChild(composeMail);

        document.body.appendChild(menu);

        const rect = e.target.getBoundingClientRect();
        menu.style.top = rect.bottom + 5 + 'px';
        menu.style.left = rect.left + 'px';

        setTimeout(() => {
            document.addEventListener('click', function closeMenu() {
                menu.remove();
                document.removeEventListener('click', closeMenu);
            });
        }, 0);
    }
    window.showSenderMenu = showSenderMenu;

    function createBody(email) {
        const container = document.createElement('div');
        container.className = 'email-body-container';

        const body = document.createElement('div');
        body.className = 'email-body';

        const hasHTML = email.Body && email.Body.trim() !== '' && email.Body !== '<pre></pre>';

        if (hasHTML && viewMode === 'html') {
            ShadowRenderer.render(body, email.Body);
        } else {
            const pre = document.createElement('pre');
            if (wordWrap) {
                pre.style.whiteSpace = 'pre-wrap';
                pre.style.wordWrap = 'break-word';
            } else {
                pre.style.whiteSpace = 'pre';
            }

            if (email.BodyText && email.BodyText.trim()) {
                pre.innerHTML = EmailUtils.linkifyText(email.BodyText);
            } else {
                pre.textContent = '[Empty Content]';
            }

            body.appendChild(pre);
        }

        container.appendChild(body);
        return container;
    }



    function switchViewMode(mode) {
        viewMode = mode;
        if (currentEmail) {
            renderEmail(currentEmail);
        }
    }

    function toggleWordWrap() {
        wordWrap = !wordWrap;
        if (currentEmail) {
            renderEmail(currentEmail);
        }
    }

    function showHeaders(email) {
        const modal = document.createElement('div');
        modal.className = 'modal-overlay';
        modal.onclick = () => modal.remove();

        const dialog = document.createElement('div');
        dialog.className = 'modal-dialog';
        dialog.onclick = (e) => e.stopPropagation();

        const title = document.createElement('h3');
        title.textContent = 'Message Headers';
        title.className = 'modal-title';

        const content = document.createElement('pre');
        content.className = 'modal-content';
        content.textContent = email.RawHeaders || 'No headers available';

        const closeBtn = document.createElement('button');
        closeBtn.textContent = 'Close';
        closeBtn.className = 'btn-close';
        closeBtn.onclick = () => modal.remove();

        dialog.appendChild(title);
        dialog.appendChild(content);
        dialog.appendChild(closeBtn);
        modal.appendChild(dialog);
        document.body.appendChild(modal);
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
});