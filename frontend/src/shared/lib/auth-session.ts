const AUTH_SESSION_KEY = 'has-session'

export function hasAuthSession() {
    return localStorage.getItem(AUTH_SESSION_KEY) === 'true'
}

export function setAuthSession() {
    localStorage.setItem(AUTH_SESSION_KEY, 'true')
}

export function clearAuthSession() {
    localStorage.removeItem(AUTH_SESSION_KEY)
}