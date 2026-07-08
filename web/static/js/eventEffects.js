// =====================================
// PARTICLES
// =====================================

const canvas = document.getElementById("particles");

if (canvas) {

    const ctx = canvas.getContext("2d");

    const particles = [];

    const mouse = {
        x: null,
        y: null,
        radius: 170
    };

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

            const light =
                document.documentElement.dataset.theme === "light" ||
                document.documentElement.dataset.bsTheme === "light";

            ctx.beginPath();

            ctx.fillStyle = light
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

        const light =
            document.documentElement.dataset.theme === "light" ||
            document.documentElement.dataset.bsTheme === "light";

        for (let a = 0; a < particles.length; a++) {

            for (let b = a + 1; b < particles.length; b++) {

                const dx = particles[a].x - particles[b].x;
                const dy = particles[a].y - particles[b].y;

                const dist = Math.sqrt(dx * dx + dy * dy);

                if (dist < 130) {

                    ctx.strokeStyle = light
                        ? `rgba(15,23,42,${0.16 - dist / 1100})`
                        : `rgba(255,255,255,${0.12 - dist / 1300})`;

                    ctx.lineWidth = 1;

                    ctx.beginPath();

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

// =====================================
// REVEAL
// =====================================

const observer = new IntersectionObserver(entries => {

    entries.forEach(entry => {

        if (entry.isIntersecting) {

            entry.target.classList.add("active");

        }

    });

}, {
    threshold: 0.15
});

document.querySelectorAll(".reveal").forEach(el => observer.observe(el));

// =====================================
// HEADER
// =====================================

const header = document.querySelector("header");

if (header) {

    header.style.background = "var(--surface)";
    header.style.backdropFilter = "blur(18px)";
    header.style.borderColor = "var(--border)";
    header.style.boxShadow = "0 12px 44px var(--shadow)";

}