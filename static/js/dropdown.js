document.addEventListener('DOMContentLoaded', function () {
    // Handle dropdown clicks
    document.querySelectorAll('.options-subitem > a').forEach(function (item) {
        item.addEventListener('click', function (e) {
            e.preventDefault();
            const parent = this.parentElement;

            if (parent.classList.contains('disabled')) {
                return;
            }

            document.querySelectorAll('.options-subitem.open').forEach(function (other) {
                if (other !== parent) {
                    other.classList.remove('open');
                }
            });

            parent.classList.toggle('open');
        });
    });

    // Close dropdowns when clicking outside
    document.addEventListener('click', function (e) {
        if (!e.target.closest('.options-subitem')) {
            document.querySelectorAll('.options-subitem.open').forEach(function (item) {
                item.classList.remove('open');
            });
        }
    });
});