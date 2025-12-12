document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchInput');
    const clearSearchIcon = document.getElementById('clearSearchIcon');
    const searchForm = document.getElementById('searchForm');

    // Очистка при клике на крестик
    if (clearSearchIcon) {
        clearSearchIcon.addEventListener('click', function() {
            searchInput.value = '';
            searchInput.focus();
            searchForm.submit();
        });
    }

});