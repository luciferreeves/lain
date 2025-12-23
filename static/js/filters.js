document.addEventListener('DOMContentLoaded', function () {
    const tagInputs = {
        'from': { input: document.getElementById('from-input'), tags: document.getElementById('from-tags'), hidden: document.getElementById('from-hidden'), values: [] },
        'to': { input: document.getElementById('to-input'), tags: document.getElementById('to-tags'), hidden: document.getElementById('to-hidden'), values: [] },
        'filename': { input: document.getElementById('filename-input'), tags: document.getElementById('filename-tags'), hidden: document.getElementById('filename-hidden'), values: [] }
    };

    const autocompleteDropdown = document.getElementById('autocomplete-dropdown');
    let activeTagInput = null;
    let autocompleteResults = [];
    let selectedIndex = -1;

    Object.keys(tagInputs).forEach(key => {
        const config = tagInputs[key];

        if (!config.input) return;

        config.input.addEventListener('keydown', function (e) {
            if (e.key === 'Enter' || e.key === ',') {
                e.preventDefault();
                const value = this.value.trim();
                if (value && !config.values.includes(value)) {
                    addTag(key, value);
                    this.value = '';
                }
            } else if (e.key === 'Backspace' && this.value === '' && config.values.length > 0) {
                removeTag(key, config.values.length - 1);
            } else if (e.key === 'ArrowDown') {
                e.preventDefault();
                navigateAutocomplete(1);
            } else if (e.key === 'ArrowUp') {
                e.preventDefault();
                navigateAutocomplete(-1);
            }
        });

        config.input.addEventListener('input', function () {
            if (this.value.length >= 2) {
                activeTagInput = key;
                showAutocomplete(this, this.value);
            } else {
                hideAutocomplete();
            }
        });

        config.input.addEventListener('blur', function () {
            setTimeout(() => hideAutocomplete(), 200);
        });
    });

    function addTag(type, value) {
        const config = tagInputs[type];
        config.values.push(value);

        const tag = document.createElement('div');
        tag.className = 'tag';
        tag.innerHTML = `
      <span>${value}</span>
      <button type="button" class="tag-remove" data-index="${config.values.length - 1}">×</button>
    `;

        tag.querySelector('.tag-remove').addEventListener('click', function () {
            removeTag(type, parseInt(this.dataset.index));
        });

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
            const tag = document.createElement('div');
            tag.className = 'tag';
            tag.innerHTML = `
        <span>${value}</span>
        <button type="button" class="tag-remove" data-index="${index}">×</button>
      `;
            tag.querySelector('.tag-remove').addEventListener('click', function () {
                removeTag(type, index);
            });
            config.tags.appendChild(tag);
        });
    }

    function showAutocomplete(input, query) {
        const suggestions = getSuggestions(query);
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

        autocompleteDropdown.innerHTML = suggestions.map((item, index) =>
            `<div class="autocomplete-item" data-index="${index}">${item}</div>`
        ).join('');

        autocompleteDropdown.querySelectorAll('.autocomplete-item').forEach(item => {
            item.addEventListener('click', function () {
                selectAutocomplete(parseInt(this.dataset.index));
            });
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

    function getSuggestions(query) {
        const mockSuggestions = [
            'john@example.com',
            'jane@example.com',
            'support@example.com',
            'admin@example.com',
            'info@example.com'
        ];
        return mockSuggestions.filter(s => s.toLowerCase().includes(query.toLowerCase()));
    }

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

    document.addEventListener('click', function (e) {
        if (!e.target.closest('.options-subitem')) {
            document.querySelectorAll('.options-subitem.open').forEach(function (item) {
                item.classList.remove('open');
            });
        }
    });

    const toggleBtn = document.getElementById('toggle-filters');
    const filters = document.getElementById('filters');
    const closeBtn = document.getElementById('close-filters');

    if (toggleBtn && filters) {
        toggleBtn.addEventListener('click', function (e) {
            e.preventDefault();
            filters.style.display = filters.style.display === 'none' ? 'block' : 'none';
        });
    }

    if (closeBtn && filters) {
        closeBtn.addEventListener('click', function () {
            filters.style.display = 'none';
        });
    }

    const clearBtn = document.getElementById('clear-filters');
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

    const datePreset = document.getElementById('date-preset');
    const customDateRange = document.getElementById('custom-date-range');
    if (datePreset && customDateRange) {
        datePreset.addEventListener('change', function () {
            customDateRange.style.display = this.value === 'custom' ? 'block' : 'none';
        });
    }

    document.addEventListener('click', function (e) {
        const filters = document.getElementById('filters');
        const toggleBtn = document.getElementById('toggle-filters');

        if (filters && toggleBtn) {
            if (!filters.contains(e.target) && !toggleBtn.contains(e.target)) {
                filters.style.display = 'none';
            }
        }
    });

    const scopeSelect = document.querySelector('select[name="scope"]');
    const customFoldersInput = document.getElementById('custom-folders-input');

    if (scopeSelect && customFoldersInput) {
        scopeSelect.addEventListener('change', function () {
            customFoldersInput.style.display = this.value === 'custom' ? 'block' : 'none';
        });
    }
});
