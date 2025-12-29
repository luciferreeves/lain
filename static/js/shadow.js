const ShadowRenderer = {
    render: function (hostElement, htmlContent) {
        let shadow = hostElement.shadowRoot;
        if (!shadow) {
            shadow = hostElement.attachShadow({ mode: 'open' });
        }

        const parser = new DOMParser();
        const doc = parser.parseFromString(htmlContent, 'text/html');

        shadow.innerHTML = '';

        if (doc.head) {
            doc.head.querySelectorAll('style').forEach(style => {
                const styleClone = style.cloneNode(true);
                let css = styleClone.textContent;
                css = css.replace(/\bbody\b/g, ':host');
                styleClone.textContent = css;
                shadow.appendChild(styleClone);
            });

            doc.head.querySelectorAll('link[rel="stylesheet"]').forEach(link => {
                shadow.appendChild(link.cloneNode(true));
            });
        }

        if (doc.body) {
            while (doc.body.firstChild) {
                shadow.appendChild(doc.body.firstChild);
            }
        }

        setTimeout(() => {
            EmailUtils.adjustContrast(shadow);
        }, 100);

        return shadow;
    }
};