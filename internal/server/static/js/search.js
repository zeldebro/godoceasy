// godoceasy - Live search & copy-to-clipboard
(function() {
    const input = document.getElementById('searchInput');
    const dropdown = document.getElementById('searchDropdown');
    let debounceTimer;

    if (!input || !dropdown) return;

    // Type badge colors
    const typeBadgeColors = {
        'Struct':    { bg: '#805ad5', text: '#fff' },
        'Interface': { bg: '#d69e2e', text: '#fff' },
        'Function':  { bg: '#38a169', text: '#fff' },
        'Type':      { bg: '#e53e3e', text: '#fff' },
        'Method':    { bg: '#3182ce', text: '#fff' },
        'Field':     { bg: '#6b7280', text: '#fff' },
        'Constant':  { bg: '#f59e0b', text: '#fff' },
        'Variable':  { bg: '#06b6d4', text: '#fff' },
        'Package':   { bg: '#326ce5', text: '#fff' }
    };

    input.addEventListener('input', function() {
        clearTimeout(debounceTimer);
        const query = this.value.trim();

        if (query.length < 2) {
            dropdown.classList.remove('active');
            dropdown.innerHTML = '';
            return;
        }

        debounceTimer = setTimeout(function() {
            fetch('/api/search?q=' + encodeURIComponent(query))
                .then(function(res) { return res.json(); })
                .then(function(results) {
                    if (!results || results.length === 0) {
                        dropdown.innerHTML = '<div class="search-dropdown-empty">No results for "<strong>' + escapeHtml(query) + '</strong>"<br><small>Try a struct, function, or type name</small></div>';
                        dropdown.classList.add('active');
                        return;
                    }

                    var count = results.length;
                    var header = '<div class="search-dropdown-header">Found ' + count + ' result' + (count > 1 ? 's' : '') + '</div>';

                    var items = results.slice(0, 12).map(function(r) {
                        var colors = typeBadgeColors[r.type] || { bg: '#718096', text: '#fff' };
                        var icon = getTypeIcon(r.type);
                        return '<a href="/pkg/' + r.package + '#' + r.anchor + '" class="search-dropdown-item">' +
                            '<div class="dropdown-item-row">' +
                            '<span class="dropdown-item-icon">' + icon + '</span>' +
                            '<span class="dropdown-item-name">' + escapeHtml(r.name) + '</span>' +
                            '<span class="dropdown-item-badge" style="background:' + colors.bg + ';color:' + colors.text + '">' + escapeHtml(r.type) + '</span>' +
                            '</div>' +
                            (r.description ? '<div class="dropdown-item-desc">' + escapeHtml(truncate(r.description, 80)) + '</div>' : '') +
                            '<div class="dropdown-item-pkg">' + escapeHtml(r.package) + '</div>' +
                            '</a>';
                    }).join('');

                    dropdown.innerHTML = header + items;
                    dropdown.classList.add('active');
                })
                .catch(function() {
                    dropdown.classList.remove('active');
                });
        }, 200);
    });

    // Close dropdown when clicking outside
    document.addEventListener('click', function(e) {
        if (!input.contains(e.target) && !dropdown.contains(e.target)) {
            dropdown.classList.remove('active');
        }
    });

    // Keyboard navigation
    input.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') {
            dropdown.classList.remove('active');
        }
        if (e.key === 'ArrowDown') {
            var first = dropdown.querySelector('.search-dropdown-item');
            if (first) { e.preventDefault(); first.focus(); }
        }
    });

    // Allow arrow key navigation in dropdown
    dropdown.addEventListener('keydown', function(e) {
        var items = dropdown.querySelectorAll('.search-dropdown-item');
        var current = document.activeElement;
        var index = Array.from(items).indexOf(current);

        if (e.key === 'ArrowDown' && index < items.length - 1) {
            e.preventDefault();
            items[index + 1].focus();
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            if (index > 0) { items[index - 1].focus(); }
            else { input.focus(); }
        } else if (e.key === 'Escape') {
            dropdown.classList.remove('active');
            input.focus();
        }
    });

    function getTypeIcon(type) {
        switch(type) {
            case 'Struct':    return '📐';
            case 'Interface': return '🔌';
            case 'Function':  return '⚡';
            case 'Type':      return '🏷️';
            case 'Method':    return '🔧';
            case 'Field':     return '📋';
            case 'Constant':  return '📌';
            case 'Variable':  return '📎';
            case 'Package':   return '📦';
            default:          return '📄';
        }
    }

    function truncate(s, max) {
        if (!s || s.length <= max) return s;
        return s.substring(0, max) + '...';
    }

    function escapeHtml(text) {
        if (!text) return '';
        var div = document.createElement('div');
        div.appendChild(document.createTextNode(text));
        return div.innerHTML;
    }
})();

// Copy to clipboard for example code blocks
function copyCode(btn) {
    var pre = btn.closest('.example-box').querySelector('pre code');
    if (!pre) return;

    var text = pre.textContent;
    navigator.clipboard.writeText(text).then(function() {
        var original = btn.textContent;
        btn.textContent = '✅ Copied!';
        btn.classList.add('copied');
        setTimeout(function() {
            btn.textContent = original;
            btn.classList.remove('copied');
        }, 2000);
    }).catch(function() {
        // Fallback for older browsers
        var textarea = document.createElement('textarea');
        textarea.value = text;
        textarea.style.position = 'fixed';
        textarea.style.opacity = '0';
        document.body.appendChild(textarea);
        textarea.select();
        try {
            document.execCommand('copy');
            var original = btn.textContent;
            btn.textContent = '✅ Copied!';
            btn.classList.add('copied');
            setTimeout(function() {
                btn.textContent = original;
                btn.classList.remove('copied');
            }, 2000);
        } catch(e) {}
        document.body.removeChild(textarea);
    });
}

// Smooth scrolling for sidebar links
document.querySelectorAll('.sidebar-list a[href^="#"]').forEach(function(link) {
    link.addEventListener('click', function(e) {
        var target = document.querySelector(this.getAttribute('href'));
        if (target) {
            e.preventDefault();
            target.scrollIntoView({ behavior: 'smooth', block: 'start' });
            history.pushState(null, null, this.getAttribute('href'));
        }
    });
});

// Highlight active sidebar item on scroll
(function() {
    var sections = document.querySelectorAll('.doc-item, .doc-section[id]');
    var sidebarLinks = document.querySelectorAll('.sidebar-sublist a');

    if (sections.length === 0 || sidebarLinks.length === 0) return;

    var ticking = false;
    window.addEventListener('scroll', function() {
        if (!ticking) {
            window.requestAnimationFrame(function() {
                var scrollPos = window.scrollY + 120;
                sections.forEach(function(section) {
                    if (section.offsetTop <= scrollPos && (section.offsetTop + section.offsetHeight) > scrollPos) {
                        sidebarLinks.forEach(function(link) {
                            link.classList.remove('sidebar-active');
                            if (link.getAttribute('href') === '#' + section.id) {
                                link.classList.add('sidebar-active');
                            }
                        });
                    }
                });
                ticking = false;
            });
            ticking = true;
        }
    });
})();

// Keyboard shortcut: "/" to focus search
document.addEventListener('keydown', function(e) {
    if (e.key === '/' && document.activeElement.tagName !== 'INPUT' && document.activeElement.tagName !== 'TEXTAREA') {
        e.preventDefault();
        var input = document.getElementById('searchInput');
        if (input) input.focus();
    }
});

// ============================================================
// Interactive Package Structure Diagram
// ============================================================
(function() {
    // Handle both package-level diagram and master overview diagram
    var containers = document.querySelectorAll('.diagram-container');
    if (!containers.length) return;

    containers.forEach(function(container) {
        var pkgName = container.getAttribute('data-package');
        if (!pkgName) return;

        if (pkgName === '__all__') {
            renderMasterDiagram(container);
        } else {
            fetch('/api/diagram/' + encodeURIComponent(pkgName))
                .then(function(res) { return res.json(); })
                .then(function(data) { renderDiagram(container, data); })
                .catch(function() {
                    container.innerHTML = '<p style="color:#888;text-align:center">Could not load diagram</p>';
                });
        }
    });

    // Master diagram — renders all packages as overview cards
    function renderMasterDiagram(el) {
        // Collect package data from the package-card elements on the page
        var cards = document.querySelectorAll('.package-card');
        if (!cards.length) {
            el.innerHTML = '<p style="color:#888;text-align:center">No packages to display</p>';
            return;
        }

        var basePath = el.getAttribute('data-base') || '';
        var html = '<div class="diagram-root">';
        html += '<div class="diagram-pkg-header">';
        html += '<span class="diagram-pkg-icon">📦</span>';
        html += '<span class="diagram-pkg-name">' + esc(basePath) + '</span>';
        html += '<div class="diagram-pkg-stats">' + cards.length + ' packages</div>';
        html += '</div>';

        html += '<div class="diagram-group">';
        html += '<div class="diagram-group-header" style="color:#1a365d">📦 All Packages (' + cards.length + ')</div>';
        html += '<div class="diagram-group-items">';

        cards.forEach(function(card) {
            var name = card.querySelector('h2');
            var stats = card.querySelectorAll('.stat');
            var href = card.getAttribute('href');
            var nameText = name ? name.textContent.trim() : 'unknown';

            html += '<a href="' + href + '" class="diagram-card" style="background:#ebf8ff;border-color:#3182ce">';
            html += '<div class="diagram-card-header" style="color:#1a365d">';
            html += '<span>📦</span> <strong>' + esc(nameText) + '</strong>';
            html += '</div>';

            if (stats.length > 0) {
                html += '<div class="diagram-card-methods">';
                stats.forEach(function(s) {
                    html += '<span class="diagram-method">' + s.textContent.trim() + '</span>';
                });
                html += '</div>';
            }

            html += '</a>';
        });

        html += '</div></div></div>';
        el.innerHTML = html;
    }

    function renderDiagram(el, data) {
        if (!data || !data.children || data.children.length === 0) {
            el.innerHTML = '<p style="color:#888;text-align:center">No items to display</p>';
            return;
        }

        var colors = {
            'struct':    { bg: '#f3e8ff', border: '#805ad5', text: '#553c9a', icon: '📐' },
            'interface': { bg: '#fefce8', border: '#d69e2e', text: '#744210', icon: '🔌' },
            'function':  { bg: '#f0fff4', border: '#38a169', text: '#22543d', icon: '⚡' },
            'type':      { bg: '#fff5f5', border: '#e53e3e', text: '#742a2a', icon: '🏷️' },
            'package':   { bg: '#ebf8ff', border: '#3182ce', text: '#1a365d', icon: '📦' }
        };

        // Group children by type
        var groups = {};
        data.children.forEach(function(child) {
            var t = child.type || 'other';
            if (!groups[t]) groups[t] = [];
            groups[t].push(child);
        });

        var html = '';

        // Package header
        html += '<div class="diagram-root">';
        html += '<div class="diagram-pkg-header">';
        html += '<span class="diagram-pkg-icon">📦</span>';
        html += '<span class="diagram-pkg-name">' + esc(data.label) + '</span>';
        if (data.count) {
            var stats = [];
            if (data.count.structs) stats.push('📐 ' + data.count.structs + ' Structs');
            if (data.count.interfaces) stats.push('🔌 ' + data.count.interfaces + ' Interfaces');
            if (data.count.functions) stats.push('⚡ ' + data.count.functions + ' Functions');
            if (data.count.types) stats.push('🏷️ ' + data.count.types + ' Types');
            if (data.count.constants) stats.push('📌 ' + data.count.constants + ' Constants');
            if (data.count.variables) stats.push('📎 ' + data.count.variables + ' Variables');
            html += '<div class="diagram-pkg-stats">' + stats.join(' &nbsp;│&nbsp; ') + '</div>';
        }
        html += '</div>';

        // Render each group
        var groupOrder = ['struct', 'interface', 'function', 'type'];
        var groupLabels = { 'struct': 'Structs', 'interface': 'Interfaces', 'function': 'Functions', 'type': 'Types' };

        groupOrder.forEach(function(gType) {
            var items = groups[gType];
            if (!items || items.length === 0) return;
            var c = colors[gType] || colors['package'];

            html += '<div class="diagram-group">';
            html += '<div class="diagram-group-header" style="color:' + c.text + '">' + c.icon + ' ' + groupLabels[gType] + ' (' + items.length + ')</div>';
            html += '<div class="diagram-group-items">';

            items.forEach(function(item) {
                var anchor = item.id;
                html += '<a href="#' + anchor + '" class="diagram-card" style="background:' + c.bg + ';border-color:' + c.border + '">';
                html += '<div class="diagram-card-header" style="color:' + c.text + '">';
                html += '<span>' + c.icon + '</span> <strong>' + esc(item.label) + '</strong>';
                html += '</div>';

                // Show fields for structs
                if (item.fields && item.fields.length > 0) {
                    html += '<div class="diagram-card-fields">';
                    var maxFields = Math.min(item.fields.length, 5);
                    for (var i = 0; i < maxFields; i++) {
                        html += '<div class="diagram-field">';
                        html += '<span class="diagram-field-name">' + esc(item.fields[i].name) + '</span>';
                        html += '<span class="diagram-field-type">' + esc(item.fields[i].type) + '</span>';
                        html += '</div>';
                    }
                    if (item.fields.length > 5) {
                        html += '<div class="diagram-field-more">+' + (item.fields.length - 5) + ' more fields</div>';
                    }
                    html += '</div>';
                }

                // Show methods for structs/interfaces
                if (item.methods && item.methods.length > 0) {
                    html += '<div class="diagram-card-methods">';
                    var maxMethods = Math.min(item.methods.length, 4);
                    for (var j = 0; j < maxMethods; j++) {
                        html += '<span class="diagram-method">.' + esc(item.methods[j]) + '()</span>';
                    }
                    if (item.methods.length > 4) {
                        html += '<span class="diagram-method-more">+' + (item.methods.length - 4) + ' more</span>';
                    }
                    html += '</div>';
                }

                // Show params/returns for functions
                if (gType === 'function') {
                    html += '<div class="diagram-card-sig">';
                    if (item.params) html += '<span class="diagram-sig-label">📥</span> ' + esc(item.params);
                    if (item.returns) html += ' <span class="diagram-sig-label">📤</span> ' + esc(item.returns);
                    html += '</div>';
                }

                html += '</a>';
            });

            html += '</div></div>';
        });

        html += '</div>';
        el.innerHTML = html;
    }

    function esc(s) {
        if (!s) return '';
        var d = document.createElement('div');
        d.appendChild(document.createTextNode(s));
        return d.innerHTML;
    }
})();

