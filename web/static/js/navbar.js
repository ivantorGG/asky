document.addEventListener('DOMContentLoaded', () => {
    showEmail()

});

async function showEmail(){
    const response = await fetch('/api/email', {
        method: "GET",
        headers: {'Content-Type': 'application/json' },
        credentials: 'include'
    })
    if (response.ok){
        const data = await response.json()
        const email = data.email


        document.getElementById('email').textContent = email
    }

    console.error('server error', response)
}