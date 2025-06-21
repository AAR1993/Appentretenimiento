(function () {
    if (localStorage.getItem('logged_in') !== 'yes') {
        window.location.replace('/login.html');
    }
})();