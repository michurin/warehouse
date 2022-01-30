// Init canvas
const canvas = document.getElementById('c');
const ctx = (() => {
  const s = 3000; // pixel size
  const l = 9; // logical size
  canvas.width = s;
  canvas.height = s;
  canvas.style.width = '700px'; // hack for retina
  const ctx = canvas.getContext('2d');
  ctx.fillStyle = '#000';
  ctx.fillRect(0, 0, s, s);
  ctx.setTransform(s / l, 0, 0, -s / l, s / 2, s / 2);
  return ctx;
})();
const scaleLimit = 0.005;

// Point
function pnt(x, y) {
  return { x, y };
}

function sum(p, q) {
  return { x: p.x + q.x, y: p.y + q.y };
}

function scale(p, k) {
  return { x: p.x * k, y: p.y * k };
}

function hypot(p) {
  return Math.hypot(p.x, p.y);
}

// Draw pentagon
function drw(o, d, color) {
  const s = hypot(d);
  ctx.save();
  ctx.fillStyle = color || `rgba(${255 * (s ** 0.6)}, 255, 0, 1)`;
  ctx.translate(o.x, o.y);
  ctx.beginPath();
  ctx.save();
  ctx.moveTo(d.x, d.y);
  for (let i = 0; i < 4; i++) {
    ctx.rotate(Math.PI * 0.4);
    ctx.lineTo(d.x, d.y);
  }
  ctx.restore();
  ctx.fill();
  if (s > 0.02) {
    ctx.fillStyle = '#fff';
    ctx.beginPath();
    ctx.arc(d.x * 0.6, d.y * 0.6, s * 0.2, 0, 2 * Math.PI, false);
    ctx.fill();
  }
  ctx.restore();
}

// Rotate clockwise and anticlockwise
const [rot, arot] = (() => {
  const a = Math.PI * 0.4;
  const c = Math.cos(a);
  const s = Math.sin(a);
  return [
    (p) => ({ x: p.x * c + p.y * s, y: p.y * c - p.x * s }),
    (p) => ({ x: p.x * c - p.y * s, y: p.y * c + p.x * s }),
  ];
})();

// Calculate scale factor
const k = (() => {
  // k * cos72 + k*k * cos36 - 1/2 = 0
  const b = Math.cos(Math.PI * 0.4);
  const a = Math.cos(Math.PI * 0.2);
  const c = -0.5;
  const d = b * b - 4 * a * c;
  const s = Math.sqrt(d);
  return (-b + s) / 2 / a;
})();

// Recursion functions

// All children, clockwise
function allChilren(o, od) {
  const r = [];
  let d = scale(od, -k);
  let s = scale(od, (1 + k)); // shift
  //  d = scale(d, -1);
  for (let i = 0; i < 5; i++) {
    r.push([sum(o, s), d]);
    s = rot(s);
    d = rot(d);
  }
  return r;
}

function rY(oo, od, leftMost, rightMost) {
  if (hypot(od) < scaleLimit) { return; }
  drw(oo, od);
  rH3(oo, od);
  const cc = allChilren(oo, od);
  // Left
  drw(...cc[2]);
  rH3(...cc[2]);
  const lcc = allChilren(...cc[2]);
  rY(...lcc[3], leftMost, false);
  if (leftMost) {
    drw(...lcc[2]);
    rH3(...lcc[2]);
    const rlcc = allChilren(...lcc[2]);
    rY(...rlcc[3], false, false);
    drw(...rlcc[2]);
    rH2(...rlcc[2]);
  } else {
    drw(...lcc[2]);
    rH2(...lcc[2]);
  }
  // Right
  drw(...cc[3]);
  rH3(...cc[3]);
  const rcc = allChilren(...cc[3]);
  rY(...rcc[2], false, rightMost);
  if (rightMost) {
    drw(...rcc[3]);
    rH3(...rcc[3]);
    const lrcc = allChilren(...rcc[3]);
    rY(...lrcc[2], false, false);
    drw(...lrcc[3]);
    rH2(...lrcc[3]);
  } else {
    drw(...rcc[3]);
    rH2(...rcc[3]);
  }
}

function rH2(o, d) {
  rH(o, rot(d));
  rH(o, arot(d));
}

function rH3(o, d) {
  rH(o, d);
  rH(o, rot(d));
  rH(o, arot(d));
}

function rH(o, d) {
  let s = scale(d, -(1 + k));
  for (; hypot(d) >= scaleLimit;) {
    d = scale(d, k * k);
    o = sum(o, s);
    drw(o, d);
    s = scale(s, k * k);
  }
}

// Central pentagon
function r0() {
  const o = pnt(0, 0);
  const d = pnt(0, 1);
  drw(o, d);
  const cc = allChilren(o, d);
  let i;
  for (i = 0; i < 5; i++) {
    rY(...cc[i], true, true);
  }
  let vd = d;
  for (i = 0; i < 5; i++) {
    rH(o, vd);
    vd = rot(vd);
  }
}

// Run
r0();
