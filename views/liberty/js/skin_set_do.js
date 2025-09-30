"use strict";

function ringo_do_regex_data(data) {
    return new RegExp('(?:^|; )' + data + '=([^;]*)');
}

function ringo_do_skin_set() {
    let cookies = document.cookie;

    if(!window.localStorage.getItem('main_css_use_sys_darkmode') || window.localStorage.getItem('main_css_use_sys_darkmode') === '1') {
        if(cookies.match(ringo_do_regex_data('main_css_darkmode')) && cookies.match(ringo_do_regex_data('main_css_darkmode'))[1] === '1') {
            if(!window.matchMedia('(prefers-color-scheme: dark)').matches) {
                document.cookie = 'main_css_darkmode=0; path=/';
                history.go(0);
            }
        }
        
        if(!cookies.match(ringo_do_regex_data('main_css_darkmode')) || cookies.match(ringo_do_regex_data('main_css_darkmode'))[1] === '0') {
            if(window.matchMedia('(prefers-color-scheme: dark)').matches) {
                document.cookie = 'main_css_darkmode=1; path=/';
                history.go(0);
            }
        }
    }
}

ringo_do_skin_set();