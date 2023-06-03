import cookie from 'js-cookie';

export function isLoggedIn(): boolean {
    const loggedIn = cookie.get('loggedIn');
    return loggedIn == "true";
}