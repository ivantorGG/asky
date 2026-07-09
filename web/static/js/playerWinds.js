// ======================================
// PARTICLES
// ======================================

const canvas = document.getElementById("particles");

if (canvas) {

    const ctx = canvas.getContext("2d");

    const particles = [];

    const mouse = {
        x: null,
        y: null,
        radius: 170
    };

    function isLightTheme() {
        return (
            document.documentElement.dataset.theme === "light" ||
            document.documentElement.dataset.bsTheme === "light"
        );
    }

    function resizeCanvas() {
        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight;
    }

    resizeCanvas();

    window.addEventListener("resize", resizeCanvas);

    window.addEventListener("mousemove", e => {
        mouse.x = e.clientX;
        mouse.y = e.clientY;
    });

    window.addEventListener("mouseleave", () => {
        mouse.x = null;
        mouse.y = null;
    });

    class Particle {

        constructor() {

            this.x = Math.random() * canvas.width;
            this.y = Math.random() * canvas.height;

            this.size = Math.random() * 2 + 1;
            this.speedX = (Math.random() - .5) * .45;
            this.speedY = (Math.random() - .5) * .45;
            this.opacity = Math.random() * .5 + .2;

        }

        update() {

            this.x += this.speedX;
            this.y += this.speedY;

            if (this.x < 0) this.x = canvas.width;
            if (this.x > canvas.width) this.x = 0;

            if (this.y < 0) this.y = canvas.height;
            if (this.y > canvas.height) this.y = 0;

            if (mouse.x !== null) {

                const dx = this.x - mouse.x;
                const dy = this.y - mouse.y;

                const dist = Math.sqrt(dx * dx + dy * dy);

                if (dist < mouse.radius) {

                    const angle = Math.atan2(dy, dx);
                    const force = (mouse.radius - dist) / mouse.radius;

                    this.x += Math.cos(angle) * force * 7;
                    this.y += Math.sin(angle) * force * 7;

                }

            }

        }

        draw() {

            ctx.beginPath();

            ctx.fillStyle = isLightTheme()
                ? `rgba(15,23,42,${this.opacity * .8})`
                : `rgba(255,255,255,${this.opacity})`;

            ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);

            ctx.fill();

        }

    }

    for (let i = 0; i < 180; i++) {
        particles.push(new Particle());
    }

    function connectParticles() {

        for (let a = 0; a < particles.length; a++) {

            for (let b = a + 1; b < particles.length; b++) {

                const dx = particles[a].x - particles[b].x;
                const dy = particles[a].y - particles[b].y;

                const dist = Math.sqrt(dx * dx + dy * dy);

                if (dist < 130) {

                    ctx.beginPath();

                    ctx.strokeStyle = isLightTheme()
                        ? `rgba(15,23,42,${0.16 - dist / 1100})`
                        : `rgba(255,255,255,${0.12 - dist / 1300})`;

                    ctx.lineWidth = 1;

                    ctx.moveTo(particles[a].x, particles[a].y);
                    ctx.lineTo(particles[b].x, particles[b].y);

                    ctx.stroke();

                }

            }

        }

    }

    function animate() {

        ctx.clearRect(0, 0, canvas.width, canvas.height);

        particles.forEach(p => {
            p.update();
            p.draw();
        });

        connectParticles();

        requestAnimationFrame(animate);

    }

    animate();

}

// ======================================
// USER NAME IN HEADER
// ======================================

async function showHeaderName() {

    try {

        const response = await fetch('/api/name', {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include'
        });

        if (!response.ok) {
            return;
        }

        const data = await response.json();
        const name = data?.name;
        const nameEl = document.getElementById('eventHeaderName');

        if (nameEl && name) {
            nameEl.textContent = name;
        }

    } catch (error) {
        console.error('Failed to load header name', error);
    }

}

async function showEventTitle() {

    try {

        const pathParts = window.location.pathname.split('/').filter(Boolean);
        const eventCode = pathParts.find(part => part.includes('-')) || pathParts[1] || '';

        if (!eventCode) {
            return;
        }

        const response = await fetch(`/api/events/${eventCode}/title`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });

        if (!response.ok) {
            return;
        }

        const data = await response.json();
        const titleEl = document.getElementById('eventHeaderTitle');

        if (titleEl && data?.title) {
            titleEl.textContent = data.title;
        }

    } catch (error) {
        console.error('Failed to load event title', error);
    }

}

showHeaderName();
showEventTitle();

// ======================================
// QUESTION CARD GLOW
// ======================================

document.addEventListener("mousemove", e => {

    document.querySelectorAll(".questionCard").forEach(card => {

        const rect = card.getBoundingClientRect();

        card.style.setProperty("--x", `${e.clientX - rect.left}px`);
        card.style.setProperty("--y", `${e.clientY - rect.top}px`);

    });

});

// ======================================
// GLASS CARD GLOW
// ======================================

document.addEventListener("mousemove", e => {

    document.querySelectorAll(".glassCard").forEach(card => {

        const rect = card.getBoundingClientRect();

        card.style.setProperty("--x", `${e.clientX - rect.left}px`);
        card.style.setProperty("--y", `${e.clientY - rect.top}px`);

    });

});

// ======================================
// AUTO RESIZE TEXTAREA
// ======================================

const textarea = document.getElementById("input2");

if (textarea) {

    textarea.addEventListener("input", () => {

        textarea.style.height = "0px";
        textarea.style.height = Math.min(textarea.scrollHeight, 180) + "px";

    });

}

// ======================================
// REVEAL
// ======================================

const observer = new IntersectionObserver(entries => {

    entries.forEach(entry => {

        if (entry.isIntersecting) {

            entry.target.classList.add("active");

        }

    });

}, {
    threshold: 0.15
});

function observeCards() {

    document.querySelectorAll(".questionCard,.glassCard").forEach(el => {

        if (!el.dataset.observed) {

            observer.observe(el);
            el.dataset.observed = "1";

        }

    });

}

observeCards();

// ======================================
// AFTER QUESTIONS RENDER
// ======================================

document.addEventListener("questionsRendered", () => {
    observeCards();
});