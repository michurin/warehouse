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

const ctx = (function() {

  // init --------------------------

  const base = 3200

  const canvas = document.getElementById('a')
  canvas.style.width = '400px'
  canvas.style.height = '400px'
  canvas.width = Math.round(base)
  canvas.height = Math.round(base)
  const context = canvas.getContext('2d')
  context.fillStyle = '#000'
  context.fillRect(0, 0, canvas.width, canvas.height)
  context.lineWidth = base / 200
  context.strokeStyle = 'rgba(255, 255, 255, 1)'
  context.lineCap = 'round'
  context.lineJoin = 'round'

  // ---- funcs ----

  const R = Math.PI / 3
  const L = -R

  function draw() {
    context.translate(step, 0)
    context.lineTo(0, 0)
  }

  function A(d) {
    d--
    if (d <= 0) {
      draw()
      return
    }
    A(d)
    context.rotate(L)
    B(d)
    context.rotate(L + L)
    B(d)
    context.rotate(R)
    A(d)
    context.rotate(R + R)
    A(d)
    A(d)
    context.rotate(R)
    B(d)
    context.rotate(L)
  }

  function B(d) {
    d--
    if (d <= 0) {
      draw()
      return
    }
    context.rotate(R)
    A(d)
    context.rotate(L)
    B(d)
    B(d)
    context.rotate(L + L)
    B(d)
    context.rotate(L)
    A(d)
    context.rotate(R + R)
    A(d)
    context.rotate(R)
    B(d)
  }

  function gosper(size, depth) {
    context.save()
    context.moveTo(0, 0)
    context.rotate(Math.atan2(Math.sqrt(3) / 2, 2.5) * (depth - 1))
    step = size / (Math.hypot(Math.sqrt(3) / 2, 2.5) ** (depth - 1))
    A(depth)
    context.stroke()
    context.restore()
  }

  // main --------------------------

  const size = base * .8
  context.translate(base / 2 - size / 2, base / 2 + size / Math.sqrt(3) / 2)

  /*
  context.save()
  context.fillStyle = 'rgba(255, 0, 0, .3)'
  context.fillRect(0, 0, size, -size / Math.sqrt(3))
  context.restore()
  */

  gosper(size, 5)

  return context
})();

(function(contextA) {

  // init --------------------------

  const base = 800
  const reflectR = 340

  const canvas = document.getElementById('b')
  canvas.style.width = '400px'
  canvas.style.height = '400px'
  canvas.width = Math.round(base)
  canvas.height = Math.round(base)
  const context = canvas.getContext('2d')
  context.fillStyle = '#030'
  context.fillRect(0, 0, canvas.width, canvas.height)

  const dataW = contextA.canvas.width
  const dataH = contextA.canvas.height
  const data = contextA.getImageData(0, 0, dataW, dataH)

  function pixel(xo, yo) {
    const x = Math.floor(dataW / 2 + xo)
    const y = Math.floor(dataH / 2 - yo)
    if (x < 0 || x >= dataW || y < 0 || y >= dataH) {
      return [0, 0, 0, 0]
    }
    const i = (x + y * dataW) * 4
    const d = data.data
    return [d[i], d[i + 1], d[i + 2], d[x + 3]]
  }

  for (let j = 0; j < base; j++) {
    for (let i = 0; i < base; i++) {
      const x = i - base / 2
      const y = j - base / 2
      const rr = (x * x + y * y) / (reflectR ** 2)
      const [r, g, b, _] = pixel(x / rr, -y / rr)
      const m = 1
      context.fillStyle = `rgb(${Math.floor(r * m)}, ${Math.floor(g * m)}, ${Math.floor(b * m)})`
      context.fillRect(i, j, 1, 1)
    }
  }
  console.log('DONE')
})(ctx)
