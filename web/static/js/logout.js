async function logout() {
    await fetch('/logout')

    location.href = '/login'
}