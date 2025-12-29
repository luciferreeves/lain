document.addEventListener('DOMContentLoaded', function () {
    document.addEventListener('click', function (e) {
        const toggleLink = e.target.closest('a');

        if (toggleLink && toggleLink.parentElement && toggleLink.parentElement.classList.contains('options-subitem')) {
            e.preventDefault();
            const parent = toggleLink.parentElement;

            if (parent.classList.contains('disabled') || parent.getAttribute('disabled') !== null) {
                return;
            }

            document.querySelectorAll('.options-subitem.open').forEach(function (other) {
                if (other !== parent) {
                    other.classList.remove('open');
                }
            });

            parent.classList.toggle('open');
            return;
        }

        // Handle clicking outside
        if (!e.target.closest('.options-subitem')) {
            document.querySelectorAll('.options-subitem.open').forEach(function (item) {
                item.classList.remove('open');
            });
        }
    });
});