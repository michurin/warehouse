(() => {
  const width = 500;
  const height = 500;

  const canvas = document.getElementsByTagName('canvas')[0];
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d');

  const scale = width / 10;

  function draw() {
    let x = Math.random() * width;
    let y = Math.random() * height;
    ctx.beginPath();
    ctx.moveTo(x, y);
    for (let i = 0; i < 50; i++) {
      const a = Math.PI * (
        Math.sin(x / scale) + Math.sin(y / scale)
      );
      x += Math.sin(a);
      y += Math.cos(a);
      ctx.lineTo(x, y);
    }
    ctx.stroke();
  }

  ctx.strokeStyle = 'rgba(100, 250, 100, .5)';
  for (let i = 0; i < 5000; i++) {
    draw();
  }
})();
