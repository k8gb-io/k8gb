(function() {
    document.addEventListener('DOMContentLoaded', () => {
        const pre = document.getElementsByTagName('pre');
        for (let i = 0; i < pre.length; i++) {
            // add the copy button iff
            // iff its parent is a div with class highlight
            if (pre[i].parentNode == null || !pre[i].parentNode.classList.contains("highlight")) {
                continue;
            }
            // and parent's parent doesn't have a class 'language-text'
            if (pre[i].parentNode.parentNode == null || pre[i].parentNode.parentNode.classList.contains("language-text")) {
                continue;
            }

            const b = document.createElement('button');
            b.className = 'clipboard';
            b.textContent = 'Copy';
            pre[i].classList.add('pre-code');
            if (pre[i].childNodes.length === 1 && pre[i].childNodes[0].nodeType === 3) {
                const code = document.createElement('code');
                code.textContent = pre[i].textContent;
                pre[i].textContent = '';
                pre[i].appendChild(code);
            }
            pre[i].appendChild(b);
        }
        new ClipboardJS('.clipboard', {
            target: (b) => {
                const p = b.parentNode;
                return p.className.includes("highlight")
                    ? p.getElementsByTagName("code")[0]
                    : p.childNodes[0];
            }
        }).on('success', (e) => {
            e.clearSelection();
            e.trigger.textContent = 'Copied';
            setTimeout(() => e.trigger.textContent = 'Copy', 2000);
        });
    });
}());
