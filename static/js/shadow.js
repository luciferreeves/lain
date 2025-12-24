/* ShadowRenderer.js - A tiny library to render safely encapsulated HTML */
const ShadowRenderer = {
    render: function (hostElement, htmlContent, options = {}) {
        // 1. Attach Shadow Root (if not exists)
        let shadow = hostElement.shadowRoot;
        if (!shadow) {
            shadow = hostElement.attachShadow({ mode: 'open' });
        }

        // 2. Parse HTML
        const parser = new DOMParser();
        const doc = parser.parseFromString(htmlContent, 'text/html');

        // 3. Rewrite Styles (Body -> Wrapper)
        const styles = doc.querySelectorAll('style');
        styles.forEach(style => {
            // Replace 'body' selector with '.mail-body-content'
            style.textContent = style.textContent.replace(/(^|[\}\s,;])body(?=[\s,\.\{])/gi, '$1.mail-body-content');
        });

        // 4. Create Wrapper
        const wrapper = document.createElement('div');
        wrapper.className = 'mail-body-content';

        // Copy attributes & move children
        if (doc.body) {
            Array.from(doc.body.attributes).forEach(attr => {
                if (attr.name === 'class') {
                    if (attr.value) wrapper.classList.add(...attr.value.split(' '));
                } else {
                    wrapper.setAttribute(attr.name, attr.value);
                }
            });
            // Handle legacy bgcolor if present
            if (doc.body.bgColor) wrapper.style.backgroundColor = doc.body.bgColor;

            while (doc.body.firstChild) wrapper.appendChild(doc.body.firstChild);
        }

        // 5. Build Shadow Content
        shadow.innerHTML = ''; // Clear previous

        // Append Head Content (Styles)
        if (doc.head) {
            while (doc.head.firstChild) shadow.appendChild(doc.head.firstChild);
        }

        // Append Body Wrapper
        shadow.appendChild(wrapper);

        // 6. Apply Default Styles
        const defaultStyle = document.createElement('style');
        defaultStyle.textContent = `
            :host { display: block; overflow: auto; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif; }
            .mail-body-content { margin: 0; padding: 16px; min-height: 100%; }
            img { max-width: 100%; height: auto; }
            a { color: #1a73e8; }
        `;
        shadow.prepend(defaultStyle);

        return shadow;
    }
};
