const ShadowRenderer = {
    render: function (hostElement, htmlContent, options = {}) {
        let shadow = hostElement.shadowRoot;
        if (!shadow) {
            shadow = hostElement.attachShadow({ mode: 'open' });
        }

        const parser = new DOMParser();
        const doc = parser.parseFromString(htmlContent, 'text/html');

        const styles = doc.querySelectorAll('style');
        styles.forEach(style => {
            style.textContent = style.textContent.replace(/(^|[\}\s,;])body(?=[\s,\.\{])/gi, '$1.mail-body-content');
        });

        const wrapper = document.createElement('div');
        wrapper.className = 'mail-body-content';

        if (doc.body) {
            Array.from(doc.body.attributes).forEach(attr => {
                if (attr.name === 'class') {
                    if (attr.value) wrapper.classList.add(...attr.value.split(' '));
                } else {
                    wrapper.setAttribute(attr.name, attr.value);
                }
            });
            if (doc.body.bgColor) wrapper.style.backgroundColor = doc.body.bgColor;

            while (doc.body.firstChild) wrapper.appendChild(doc.body.firstChild);
        }
        shadow.innerHTML = '';

        if (doc.head) {
            while (doc.head.firstChild) shadow.appendChild(doc.head.firstChild);
        }
        shadow.appendChild(wrapper);
        const defaultStyle = document.createElement('style');
        defaultStyle.textContent = `
            :host { display: block; overflow: auto; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif; }
            img { max-width: 100%; height: auto; }
            a { color: #1a73e8; }
        `;
        shadow.prepend(defaultStyle);

        return shadow;
    }
};
