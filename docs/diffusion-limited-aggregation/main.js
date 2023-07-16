// ------- geometry -------

function rVector(r) {
  var a = Math.PI * 2 * Math.random();
  return {
    x: r * Math.cos(a),
    y: r * Math.sin(a),
  };
}

function add(v, w) {
  return {
    x: v.x + w.x,
    y: v.y + w.y,
  };
}

function sub(v, w) {
  return {
    x: v.x - w.x,
    y: v.y - w.y,
  };
}

function hypot(p) {
  return Math.hypot(p.x, p.y);
}

// ------- storage -------

function store(s, p) {
  const a = Math.floor(p.x / 5);
  const b = Math.floor(p.y / 5);
  for (var i = -1; i < 2; i++) {
    for (var j = -1; j < 2; j++) {
      var k = `${a + i}x${b + j}`;
      if (s[k] === undefined) {
        s[k] = [];
      }
      s[k].push(p);
    }
  }
}

function lookup(s, p) {
  const a = Math.floor(p.x / 5);
  const b = Math.floor(p.y / 5);
  const k = `${a}x${b}`;
  const qq = s[k];
  if (qq === undefined) {
    return undefined;
  }
  for (j = 0; j < qq.length; j++) {
    var q = qq[j];
    if (hypot(sub(q, p)) < 2) {
      return q;
    }
  }
  return undefined;
}

// ------- main -------

const size = 500;

const canvas = document.getElementById("a");
canvas.width = size;
canvas.height = size;
const context = canvas.getContext("2d");
context.translate(size / 2, size / 2); // FORCE 200x200 with (0, 0) in the middle
context.scale(size / 200, size / 200);
context.lineWidth = 2;
context.strokeStyle = '#08f'; // will be redefined later
context.lineCap = 'round';

const points = {};
store(points, { x: 0, y: 0 });
var radius = 4;
var giveUpRadius = 6;

function put(iteration) {
  var p = rVector(radius);
  //console.log('NEW', p, radius);
  while (hypot(p) < giveUpRadius) {
    p = add(p, rVector(1));
    q = lookup(points, p);
    if (q) {
      radius = Math.max(radius, hypot(p));
      giveUpRadius = radius * 1.1;
      store(points, p);
      context.beginPath();
      context.moveTo(q.x, q.y);
      context.lineTo(p.x, p.y);
      context.strokeStyle = `rgb(${Math.floor((1 - iteration / 5000) * 255)}, 255, 0)`;
      context.stroke();
      //console.log('STICK', p, q, hypot(p));
      break;
    }
  }
  iteration--;
  if (iteration > 0) {
    setTimeout(function() { put(iteration); }, 0);
  }
}

put(5000)
