const EmailUtils = {
    getProfilePicture: async function (email, name) {
        const gravatarUrl = await this.checkGravatar(email);
        if (gravatarUrl) return gravatarUrl;

        let domain = email.split('@')[1];

        // Remove all subdomains from the domain
        const domainParts = domain.split('.');
        if (domainParts.length > 2) {
            domainParts.shift();
            domain = domainParts.join('.');
        }

        return `https://t2.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=http://${domain}&size=128`;
    },

    checkGravatar: async function (email) {
        const hash = await this.sha256(email.toLowerCase().trim());
        const testUrl = `https://www.gravatar.com/avatar/${hash}?s=128&d=404`;

        return new Promise((resolve) => {
            const img = new Image();
            img.onload = () => resolve(testUrl);
            img.onerror = () => resolve(null);
            img.src = testUrl;
        });
    },

    sha256: async function (message) {
        const msgBuffer = new TextEncoder().encode(message);
        const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
        const hashArray = Array.from(new Uint8Array(hashBuffer));
        return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    },

    getInitials: function (name, email) {
        if (name && name !== email) {
            const parts = name.trim().split(' ');
            if (parts.length >= 2) {
                return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
            }
            return name.substring(0, 2).toUpperCase();
        }
        return email.substring(0, 2).toUpperCase();
    },

    formatEmailDisplay: function (name, email, showAddress) {
        if (!name || name === email) {
            return email;
        }

        if (showAddress) {
            return `${name} <${email}>`;
        }

        return name;
    },

    linkifyText: function (text) {
        const urlRegex = /(https?:\/\/[^\s]+)/g;
        const emailRegex = /([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)/g;

        return text
            .replace(urlRegex, '<a href="$1" target="_blank" style="color: var(--accent-primary); text-decoration: underline;">$1</a>')
            .replace(emailRegex, '<a href="mailto:$1" style="color: var(--accent-primary); text-decoration: underline;">$1</a>');
    },

    adjustContrast: function (container) {
        const allElements = container.querySelectorAll('*');

        allElements.forEach(el => {
            const style = window.getComputedStyle(el);
            const color = style.color;
            const bgColor = style.backgroundColor;

            const textLum = this.getLuminance(color);
            const bgLum = this.getLuminance(bgColor);

            const isTransparent = bgColor === 'rgba(0, 0, 0, 0)' || bgColor === 'transparent';

            let actualBgLum = bgLum;
            if (isTransparent) {
                let parent = el.parentElement;
                while (parent && (window.getComputedStyle(parent).backgroundColor === 'rgba(0, 0, 0, 0)' ||
                    window.getComputedStyle(parent).backgroundColor === 'transparent')) {
                    parent = parent.parentElement;
                }
                if (parent) {
                    actualBgLum = this.getLuminance(window.getComputedStyle(parent).backgroundColor);
                }
            }

            if (actualBgLum > 0.8 && textLum > 0.55) {
                el.style.setProperty('color', '#000000', 'important');
            } else if (actualBgLum < 0.3 && textLum < 0.45) {
                el.style.setProperty('color', '#e8e8f0', 'important');
            }
        });
    },

    getLuminance: function (color) {
        const rgb = color.match(/\d+/g);
        if (!rgb || rgb.length < 3) return 1;

        const [r, g, b] = rgb.map(val => {
            const v = val / 255;
            return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4);
        });

        return 0.2126 * r + 0.7152 * g + 0.0722 * b;
    }
};