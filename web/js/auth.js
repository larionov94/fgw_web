// HTML5 placeholder support for older browsers
(function() {
    if (!('placeholder' in document.createElement('input'))) {
        let inputs = document.querySelectorAll('input[placeholder]');
        for (var i = 0; i < inputs.length; i++) {
            var input = inputs[i];
            var placeholder = input.getAttribute('placeholder');

            input.value = placeholder;
            input.style.color = '#999';

            input.addEventListener('focus', function() {
                if (this.value === this.getAttribute('placeholder')) {
                    this.value = '';
                    this.style.color = '';
                }
            });

            input.addEventListener('blur', function() {
                if (this.value === '') {
                    this.value = this.getAttribute('placeholder');
                    this.style.color = '#999';
                }
            });
        }
    }
})();