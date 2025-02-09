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

// How to check
// eslint --no-eslintrc --fix main.js

const unitSize = 20
const fieldWidthUnits = 25
const fieldHeightUnits = 25

const creaturesCounts = [400, 400]

// ----- init

const interactionMatrics = (() => {
  return [ // + — attract
    [{ r1: .3, r2: .6, m: .1 }, { r1: .4, r2: .6, m: .06 }],
    [{ r1: .4, r2: .6, m: -.02 }, { r1: .3, r2: .6, m: .1 }],
  ]
  const mtx = [] // random. Is not used right now
  for (let i = 0; i < creaturesCounts.length; i++) {
    const m = []
    for (let j = 0; j < creaturesCounts.length; j++) {
      const a = 1 - 2 * Math.random()
      let p = .2 // Math.random()
      let q = .5 // Math.random()
      if (p > q) {
        t = p
        p = q
        q = t
      }
      m.push({
        r1: p,
        r2: q,
        m: a,
      })
    }
    mtx.push(m)
  }
  console.log(mtx)
  return mtx
})()

const canvas = document.getElementById('a')
const context = canvas.getContext('2d')

canvas.width = Math.round(unitSize * fieldWidthUnits)
canvas.height = Math.round(unitSize * fieldHeightUnits)

const creatures = (() => { // It should not be const in future?
  const c = []
  for (let tp = 0; tp < creaturesCounts.length; tp++) {
    const clr = `hsl(${Math.round(tp * 360 / creaturesCounts.length)} 100% 50%)`
    for (let i = 0; i < creaturesCounts[tp]; i++) {
      const v = Math.random()
      const d = Math.random() * Math.PI * 2
      c.push({
        x: fieldWidthUnits * Math.random(),
        y: fieldHeightUnits * Math.random(),
        vx: v * Math.cos(d),
        vy: v * Math.sin(d),
        tp: tp,
        clr: clr,
      })
    }
  }
  return c
})()

function draw() {
  context.clearRect(0, 0, unitSize * fieldWidthUnits, unitSize * fieldHeightUnits)
  // IDEA:
  // context.filter = 'brightness(80%)'
  // context.drawImage(canvas, 0, 0)
  // context.filter = ''
  for (let i = 0; i < creatures.length; i++) {
    const e = creatures[i]
    context.fillStyle = e.clr
    context.fillRect(e.x * unitSize - 1, e.y * unitSize - 1, 3, 3)
  }
}

function move(ms) {
  for (let i = 0; i < creatures.length; i++) {
    const e = creatures[i]
    e.x = (e.x + ms * e.vx / 1000 + fieldWidthUnits) % fieldWidthUnits
    e.y = (e.y + ms * e.vy / 1000 + fieldHeightUnits) % fieldHeightUnits
  }
}

function accelerate(ms) {
  const dt = ms * .001
  for (let i = 0; i < creatures.length; i++) {
    const e = creatures[i]
    let fx = 0
    let fy = 0
    for (let j = 0; j < creatures.length; j++) {
      const g = creatures[j]
      const dx = e.x - g.x
      const dy = e.y - g.y
      const r = Math.hypot(dx, dy)
      if (i !== j && r < 1) {
        const int = interactionMatrics[e.tp][g.tp]
        let f = 0
        if (r < int.r1) {
          f = -1 + r / int.r1
        } else if (r < int.r2) {
          f = int.m * (r - int.r1) / (int.r2 - int.r1) // TODO: pre calculate denominator
        } else {
          f = int.m * (1 - r) / (1 - int.r2)
        }
        f = -f // TODO: dirty hack
        fx += f * dx / r
        fy += f * dy / r
      }
    }
    fx -= e.vx * .4
    fy -= e.vy * .4
    e.vx += dt * fx
    e.vy += dt * fy
  }
}

let prevMs = 0
function animate(ms) {
  if (prevMs === 0) {
    prevMs = ms
  }
  accelerate(ms - prevMs)
  move(ms - prevMs)
  draw()
  prevMs = ms
  window.requestAnimationFrame(animate)
}

window.requestAnimationFrame(animate)

