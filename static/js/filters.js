document.addEventListener('DOMContentLoaded', function () {
    const tagInputs = {
        'from': {
            input: document.getElementById('from-input'),
            tags: document.getElementById('from-tags'),
            hidden: document.getElementById('from-hidden'),
            values: []
        },
        'to': {
            input: document.getElementById('to-input'),
            tags: document.getElementById('to-tags'),
            hidden: document.getElementById('to-hidden'),
            values: []
        },
        'filename': {
            input: document.getElementById('filename-input'),
            tags: document.getElementById('filename-tags'),
            hidden: document.getElementById('filename-hidden'),
            values: []
        }
    };

    const autocompleteDropdown = document.getElementById('autocomplete-dropdown');
    let activeTagInput = null;
    let autocompleteResults = [];
    let selectedIndex = -1;

    // Initialize tag inputs
    Object.keys(tagInputs).forEach(key => {
        const config = tagInputs[key];
        if (!config.input) return;

        config.input.addEventListener('keydown', handleTagInput.bind(null, key));
        config.input.addEventListener('input', handleAutocomplete.bind(null, key));
        config.input.addEventListener('blur', () => setTimeout(hideAutocomplete, 200));
    });

    function handleTagInput(type, e) {
        const config = tagInputs[type];

        if (e.key === 'Enter' || e.key === ',') {
            e.preventDefault();
            const value = e.target.value.trim();
            if (value && !config.values.includes(value)) {
                addTag(type, value);
                e.target.value = '';
            }
        } else if (e.key === 'Backspace' && e.target.value === '' && config.values.length > 0) {
            removeTag(type, config.values.length - 1);
        } else if (e.key === 'ArrowDown') {
            e.preventDefault();
            navigateAutocomplete(1);
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            navigateAutocomplete(-1);
        }
    }

    function handleAutocomplete(type, e) {
        if (e.target.value.length >= 2) {
            activeTagInput = type;
            showAutocomplete(e.target, e.target.value);
        } else {
            hideAutocomplete();
        }
    }

    function addTag(type, value) {
        const config = tagInputs[type];
        config.values.push(value);

        const tag = document.createElement('div');
        tag.className = 'tag';

        const span = document.createElement('span');
        span.textContent = value;

        const button = document.createElement('button');
        button.type = 'button';
        button.className = 'tag-remove';
        button.dataset.index = config.values.length - 1;
        button.textContent = '×';
        button.addEventListener('click', function () {
            removeTag(type, parseInt(this.dataset.index));
        });

        tag.appendChild(span);
        tag.appendChild(button);
        config.tags.appendChild(tag);
        config.hidden.value = config.values.join(',');
    }

    function removeTag(type, index) {
        const config = tagInputs[type];
        config.values.splice(index, 1);
        renderTags(type);
        config.hidden.value = config.values.join(',');
    }

    function renderTags(type) {
        const config = tagInputs[type];
        config.tags.innerHTML = '';
        config.values.forEach((value, index) => {
            addTag(type, value);
        });
    }

    function showAutocomplete(input, query) {
        // TODO: Replace with actual API call to fetch contacts
        const suggestions = [];

        if (suggestions.length === 0) {
            hideAutocomplete();
            return;
        }

        autocompleteResults = suggestions;
        selectedIndex = -1;

        const rect = input.getBoundingClientRect();
        autocompleteDropdown.style.top = (rect.bottom + window.scrollY) + 'px';
        autocompleteDropdown.style.left = rect.left + 'px';
        autocompleteDropdown.style.width = rect.width + 'px';
        autocompleteDropdown.innerHTML = '';

        suggestions.forEach((item, index) => {
            const div = document.createElement('div');
            div.className = 'autocomplete-item';
            div.dataset.index = index;
            div.textContent = item;
            div.addEventListener('click', () => selectAutocomplete(index));
            autocompleteDropdown.appendChild(div);
        });

        autocompleteDropdown.style.display = 'block';
    }

    function hideAutocomplete() {
        if (autocompleteDropdown) {
            autocompleteDropdown.style.display = 'none';
        }
        activeTagInput = null;
        autocompleteResults = [];
        selectedIndex = -1;
    }

    function navigateAutocomplete(direction) {
        if (autocompleteResults.length === 0) return;

        selectedIndex = Math.max(-1, Math.min(autocompleteResults.length - 1, selectedIndex + direction));

        const items = autocompleteDropdown.querySelectorAll('.autocomplete-item');
        items.forEach((item, index) => {
            item.classList.toggle('active', index === selectedIndex);
        });

        if (selectedIndex >= 0) {
            items[selectedIndex].scrollIntoView({ block: 'nearest' });
        }
    }

    function selectAutocomplete(index) {
        if (activeTagInput && autocompleteResults[index]) {
            addTag(activeTagInput, autocompleteResults[index]);
            tagInputs[activeTagInput].input.value = '';
            hideAutocomplete();
            tagInputs[activeTagInput].input.focus();
        }
    }

    // Filter controls
    const toggleBtn = document.getElementById('toggle-filters');
    const filters = document.getElementById('filters');
    const closeBtn = document.getElementById('close-filters');
    const clearBtn = document.getElementById('clear-filters');

    if (toggleBtn && filters) {
        toggleBtn.addEventListener('click', function (e) {
            e.preventDefault();
            filters.style.display = filters.style.display === 'none' ? 'block' : 'none';
        });
    }

    if (closeBtn && filters) {
        closeBtn.addEventListener('click', () => filters.style.display = 'none');
    }

    if (clearBtn) {
        clearBtn.addEventListener('click', function () {
            Object.keys(tagInputs).forEach(key => {
                tagInputs[key].values = [];
                tagInputs[key].tags.innerHTML = '';
                tagInputs[key].hidden.value = '';
                tagInputs[key].input.value = '';
            });
            document.querySelectorAll('.filter-input').forEach(input => input.value = '');
            document.querySelectorAll('.flag-option input[type="checkbox"]').forEach(cb => cb.checked = false);
        });
    }

    // Date preset handling
    const datePreset = document.getElementById('date-preset');
    const customDateRange = document.getElementById('custom-date-range');

    if (datePreset && customDateRange) {
        datePreset.addEventListener('change', function () {
            customDateRange.style.display = this.value === 'custom' ? 'block' : 'none';
        });
    }

    // Scope handling
    const scopeSelect = document.querySelector('select[name="scope"]');
    const customFoldersInput = document.getElementById('custom-folders-input');

    if (scopeSelect && customFoldersInput) {
        scopeSelect.addEventListener('change', function () {
            customFoldersInput.style.display = this.value === 'custom' ? 'block' : 'none';
        });
    }

    // Close filters when clicking outside
    document.addEventListener('click', function (e) {
        if (filters && toggleBtn) {
            if (!filters.contains(e.target) && !toggleBtn.contains(e.target)) {
                filters.style.display = 'none';
            }
        }
    });
});