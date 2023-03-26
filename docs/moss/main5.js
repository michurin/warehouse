(function() {
  const width = 500;
  const height = 500;

  const canvas = document.getElementsByTagName('canvas')[4];
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d');

  function fourierFactory() {
    const cc = [];
    for (let i = 0; i < 6; i++) {
      const r = (Math.random() + 0.5) / width * Math.PI * 2 * 2;
      const a = Math.random() * 2 * Math.PI;
      cc.push({
        x: r * Math.cos(a),
        y: r * Math.sin(a),
        phase: Math.random() * 1e5,
      });
    }
    return (x, y) => {
      let s = 0;
      cc.forEach((c) => {
        s += Math.sin(x * c.x + y * c.y + c.phase);
      });
      return s / cc.length;
    };
  }

  let fourier;

  function draw(rot) {
    let x = Math.random() * width;
    let y = Math.random() * height;
    ctx.beginPath();
    ctx.moveTo(x, y);
    for (let i = 0; i < 50; i++) {
      const a = 4 * Math.PI * fourier(x, y) + rot;
      x += Math.sin(a);
      y += Math.cos(a);
      ctx.lineTo(x, y);
    }
    ctx.stroke();
  }

  function fix() { return new Promise((a, _) => { setTimeout(() => a(), 0); }); }

  async function drawFigure() {
    fourier = fourierFactory();
    ctx.clearRect(0, 0, width, height);
    ctx.lineWidth = 1;
    for (let i = 0; i < 5000; i++) {
      ctx.strokeStyle = 'rgba(100, 250, 100, .5)';
      draw(0);
      ctx.strokeStyle = 'rgba(250, 250, 100, .5)';
      draw(Math.PI * 2 / 3);
      ctx.strokeStyle = 'rgba(100, 250, 250, .5)';
      draw(Math.PI * 4 / 3);
      if (i % 100 === 0) { await fix(); }
    }
  }

  drawFigure();
  canvas.onclick = drawFigure;
}());

