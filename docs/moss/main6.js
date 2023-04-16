(function() {
  const width = 500;
  const height = 500;

  const canvas = document.getElementsByTagName('canvas')[5];
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d');

  const pivotPoints = [];

  function draw(mx) {
    let x = Math.random() * width;
    let y = Math.random() * height;
    ctx.beginPath();
    ctx.moveTo(x, y);
    let vx = 0;
    let vy = 0;
    for (let j = 0; j < 50; j++) {
      let fx = 0;
      let fy = 0;
      for (let i = 0; i < pivotPoints.length; i++) {
        let rx = x - pivotPoints[i].x;
        let ry = y - pivotPoints[i].y;
        [rx, ry] = [rx * mx.xx + ry * mx.xy, rx * mx.yx + ry * mx.yy];
        let r = Math.hypot(rx, ry);
        let r3 = r * r * r * r;
        rx /= r3;
        ry /= r3;
        fx += rx;
        fy += ry;
      }
      let f = Math.hypot(fx, fy);
      vx += fx / f;
      vy += fy / f;
      let v = Math.hypot(vx, vy);
      x += vx / v;
      y += vy / v;
      ctx.lineTo(x, y);
    }
    ctx.stroke();
  }

  function fix() { return new Promise((a, _) => { setTimeout(() => a(), 0); }); }

  function matrix(a) {
    return {
      xx: Math.cos(a),
      xy: -Math.sin(a),
      yx: Math.sin(a),
      yy: Math.cos(a),
    };
  }

  async function drawFigure() {
    pivotPoints.length = 0;
    for (let i = 0; i < 32; i++) {
      pivotPoints.push({
        x: Math.random() * width,
        y: Math.random() * height,
      });
    }
    const ma = matrix(Math.PI / 2 * .6);
    const mb = matrix(Math.PI / 2 * 1.6); // perpendicular
    ctx.clearRect(0, 0, width, height);
    ctx.lineWidth = 1;
    for (let i = 0; i < 5000; i++) {
      ctx.strokeStyle = 'rgba(100, 250, 100, .5)';
      draw(ma);
      ctx.strokeStyle = 'rgba(250, 250, 100, .5)';
      draw(mb);
      if (i % 100 === 0) { await fix(); }
    }
  }

  drawFigure();
  canvas.onclick = drawFigure;
}());
