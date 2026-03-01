export default class ParticleNetwork {
  constructor(canvas, options = {}) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d');
    this.particles = [];
    this.count = options.count || 80;
    this.maxDistance = options.maxDistance || 150;
    this.speed = options.speed || 0.3;
    this.mouseX = -1000;
    this.mouseY = -1000;
    this.animFrame = null;
    this.init();
    this.bindMouse();
  }

  init() {
    this.resize();
    window.addEventListener('resize', () => this.resize());
    this._createParticles();
  }

  _createParticles() {
    this.particles = [];
    const tealCount = Math.floor(this.count * 0.2);
    const regularCount = this.count - tealCount;

    for (let i = 0; i < regularCount; i++) {
      this.particles.push(this._makeParticle(false));
    }
    for (let i = 0; i < tealCount; i++) {
      this.particles.push(this._makeParticle(true));
    }
  }

  _makeParticle(isTeal) {
    const speed = this.speed;
    return {
      x: Math.random() * this.canvas.width,
      y: Math.random() * this.canvas.height,
      vx: (Math.random() - 0.5) * speed * 2,
      vy: (Math.random() - 0.5) * speed * 2,
      radius: 1.5 + Math.random() * 1.5,
      isTeal,
      color: isTeal ? 'rgba(0,212,170,0.8)' : 'rgba(148,163,184,0.4)',
      opacity: isTeal ? 0.8 : 0.4,
    };
  }

  resize() {
    const w = this.canvas.offsetWidth;
    const h = this.canvas.offsetHeight;
    if (w > 0 && h > 0) {
      this.canvas.width = w;
      this.canvas.height = h;
    }
  }

  bindMouse() {
    document.addEventListener('mousemove', (e) => {
      const rect = this.canvas.getBoundingClientRect();
      this.mouseX = e.clientX - rect.left;
      this.mouseY = e.clientY - rect.top;
    });
  }

  draw() {
    const ctx = this.ctx;
    const w = this.canvas.width;
    const h = this.canvas.height;

    ctx.clearRect(0, 0, w, h);

    // Draw connections between nearby particles
    for (let i = 0; i < this.particles.length; i++) {
      for (let j = i + 1; j < this.particles.length; j++) {
        const a = this.particles[i];
        const b = this.particles[j];
        const dx = a.x - b.x;
        const dy = a.y - b.y;
        const dist = Math.sqrt(dx * dx + dy * dy);

        if (dist < this.maxDistance) {
          const alpha = (1 - dist / this.maxDistance) * 0.35;
          const isTealLine = a.isTeal || b.isTeal;
          if (isTealLine) {
            ctx.strokeStyle = `rgba(0,212,170,${alpha})`;
          } else {
            ctx.strokeStyle = `rgba(148,163,184,${alpha * 0.7})`;
          }
          ctx.lineWidth = isTealLine ? 0.8 : 0.5;
          ctx.beginPath();
          ctx.moveTo(a.x, a.y);
          ctx.lineTo(b.x, b.y);
          ctx.stroke();
        }
      }
    }

    // Mouse repulsion / attraction effect
    for (let i = 0; i < this.particles.length; i++) {
      const p = this.particles[i];
      const dx = p.x - this.mouseX;
      const dy = p.y - this.mouseY;
      const dist = Math.sqrt(dx * dx + dy * dy);
      const mouseRadius = 120;

      if (dist < mouseRadius && dist > 0) {
        const force = (mouseRadius - dist) / mouseRadius;
        const ax = (dx / dist) * force * 0.4;
        const ay = (dy / dist) * force * 0.4;
        p.vx += ax;
        p.vy += ay;
      }

      // Dampen velocity to prevent runaway
      const maxSpeed = this.speed * 3;
      const currentSpeed = Math.sqrt(p.vx * p.vx + p.vy * p.vy);
      if (currentSpeed > maxSpeed) {
        p.vx = (p.vx / currentSpeed) * maxSpeed;
        p.vy = (p.vy / currentSpeed) * maxSpeed;
      }

      // Move particle
      p.x += p.vx;
      p.y += p.vy;

      // Wrap around edges
      if (p.x < -p.radius) p.x = w + p.radius;
      else if (p.x > w + p.radius) p.x = -p.radius;
      if (p.y < -p.radius) p.y = h + p.radius;
      else if (p.y > h + p.radius) p.y = -p.radius;

      // Draw particle
      ctx.beginPath();
      ctx.arc(p.x, p.y, p.radius, 0, Math.PI * 2);
      ctx.fillStyle = p.color;
      ctx.fill();

      // Glow for teal particles
      if (p.isTeal) {
        ctx.beginPath();
        ctx.arc(p.x, p.y, p.radius * 2.5, 0, Math.PI * 2);
        const grd = ctx.createRadialGradient(p.x, p.y, 0, p.x, p.y, p.radius * 2.5);
        grd.addColorStop(0, 'rgba(0,212,170,0.15)');
        grd.addColorStop(1, 'rgba(0,212,170,0)');
        ctx.fillStyle = grd;
        ctx.fill();
      }
    }

    this.animFrame = requestAnimationFrame(() => this.draw());
  }

  start() {
    this.draw();
  }

  stop() {
    if (this.animFrame) {
      cancelAnimationFrame(this.animFrame);
      this.animFrame = null;
    }
  }
}
