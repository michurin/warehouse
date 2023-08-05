/*eslint-env es2017*/
/*eslint indent: ["error", 2]*/
/*eslint eqeqeq: ["error", "always"]*/
/*eslint prefer-const: "error"*/
/*eslint no-var: "error"*/
/*eslint one-var: ["error", "never"]*/
/*eslint semi: ["error", "never"]*/
/*eslint quotes: ["error", "single"]*/
/*eslint prefer-arrow-callback: "error"*/
/*eslint arrow-body-style: "error"*/

// init --------------------------

const base = parseInt((new URLSearchParams(window.location.search)).get('s') / 4.5) || 256

const canvas = document.getElementById('a')
canvas.width = Math.round(base * 4.5)
canvas.height = Math.round(base * 3)
const context = canvas.getContext('2d')
context.fillStyle = '#000'
context.fillRect(0, 0, canvas.width, canvas.height)
context.lineWidth = 1
context.strokeStyle = '#08f' // will be redefined later
context.lineCap = 'round'
context.lineJoin = 'round'

// funcs -------------------------

// geometry

function tmod(t) {
  return Math.hypot(t[2], t[3])
}

function trans(t, v) { // transform vector
  const [sx, sy, xx, xy] = t
  const [x, y] = v
  const [yx, yy] = [-xy, xx]
  return [
    sx + x * xx + y * yx,
    sy + x * xy + y * yy,
  ]
}

function transt(t, x) { // transform transformation
  const [sx, sy, xx, xy] = x
  const [ax, ay] = trans(t, [sx, sy])
  const [bx, by] = trans([0, 0, t[2], t[3]], [xx, xy]) // just rotate, don't shift
  return [ax, ay, bx, by]
}

// graphic

function triangle(t) {
  const grd = context.createLinearGradient(...trans(t, [1, 0]), ...trans(t, [0, 0]))
  grd.addColorStop(0, '#00f')
  grd.addColorStop(.2, '#008')
  grd.addColorStop(1, '#000')
  context.beginPath()
  context.moveTo(...trans(t, [0, 0]))
  context.lineTo(...trans(t, [1, 0]))
  context.lineTo(...trans(t, [0, -1]))
  context.fillStyle = grd
  context.closePath()
  context.fill()
  context.stroke()
}

function vector(t) { // for debugging; you may want to turn off filling in triangle() to see vectors
  const p = .1
  context.save()
  context.beginPath()
  context.moveTo(...trans(t, [0, 0]))
  context.lineTo(...trans(t, [1, 0]))
  context.lineTo(...trans(t, [1 - p, p]))
  context.lineTo(...trans(t, [1 - p, -p]))
  context.lineTo(...trans(t, [1, 0]))
  context.lineWidth = 2
  context.strokeStyle = '#f00'
  context.stroke()
  context.restore()
}

// fractal

function tooth(t) {
  triangle(t)
  spiral(transt(t, [0, 0, 0, .5]))
  factory(transt(t, [0, -1, 0, 1]))
  factory(transt(t, [.5, 0, .5, 0]))
}

function factory(t) {
  if (tmod(t) < 1) { return }
  tooth(transt(t, [.5, .25, .25, -.25]))
  factory(transt(t, [0, 0, .25, 0]))
  factory(transt(t, [.75, 0, .25, 0]))
}

function spiral(t) {
  if (tmod(t) < 1) { return }
  triangle(t)
  spiral(transt(t, [1, 0, .5, -.5]))
  triangle(transt(t, [.5, .5, .5, -.5]))
  spiral(transt(t, [.5, .5, .25, .25]))
  factory(transt(t, [.75, .25, .25, -.25]))
  factory(transt(t, [0, 0, .5, .5]))
  factory(transt(t, [.5, -.5, -.5, -.5]))
}

function core(t) {
  if (tmod(t) < 1) { return }
  triangle(t)
  factory(transt(t, [1, 0, -1, -1]))
  spiral(transt(t, [0, 0, 0, .5]))
  factory(transt(t, [0, -1, 0, 1]))
  factory(transt(t, [.5, 0, .5, 0]))
  core(transt(t, [1.5, .5, .5, -.5]))
}

function seed(t) {
  core(t)
  core(transt(t, [0, -1.5, -.5, 0]))
}

// main --------------------------

seed([base * 3, base, -base, 0])
