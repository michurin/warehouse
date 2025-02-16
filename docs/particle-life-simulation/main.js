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

const partitionSize = 20

const arenaWidth = 40
const arenaHeight = 40

function createEmptyArena() {
  const a = []
  for (let i = 0; i < arenaHeight; i++) {
    const r = []
    for (let j = 0; j < arenaWidth; j++) {
      r.push([])
    }
    a.push(r)
  }
  return a
}

let arena = createEmptyArena()

// ---- fill arena

for (let tp = 0; tp < 4; tp++) {
  for (let s = 0; s < 500; s++) {
    const x = arenaWidth * Math.random()
    const y = arenaHeight * Math.random()
    console.log(x, y)
    const v = 1 // partitions/second
    const d = Math.random() * Math.PI * 2
    arena[Math.floor(y)][Math.floor(x)].push({
      x: x % 1, // coordinates, related partition
      y: y % 1,
      vx: v * Math.cos(d),
      vy: v * Math.sin(d),
      tp: tp,
      clr: `hsl(${Math.round(tp * 360 / 4)} 100% 50%)`
    })
  }
}
console.log(arena)

// ---- setup canvas

const canvas = document.getElementById('a')
const context = canvas.getContext('2d')

canvas.width = Math.round(partitionSize * arenaWidth)
canvas.height = Math.round(partitionSize * arenaHeight)

function drawArena() {
  context.clearRect(0, 0, partitionSize * arenaWidth, partitionSize * arenaHeight)
  for (let pr = 0; pr < arenaHeight; pr++) { // loop over partitions rows
    for (let p = 0; p < arenaWidth; p++) { // loop over partitions
      const pv = arena[pr][p]
      for (let i = 0; i < pv.length; i++) { // loop over members of partition
        const v = pv[i]
        context.fillStyle = v.clr
        context.fillRect((p + v.x) * partitionSize - 1, (pr + v.y) * partitionSize - 1, 3, 3)
      }
    }
  }
}

const r1 = .2
const r2 = .7
const forceMatrix = [
  [.4, .2, 0, 0],
  [-.05, .4, 0, 0],
  [0, 0, .4, .2],
  [0, 0, -.005, .4],
]

function accelerate(dt) {
  for (let tpr = 0; tpr < arenaHeight; tpr++) { // loop over target partition rows
    for (let tp = 0; tp < arenaWidth; tp++) { // loop over target partitions
      const tpv = arena[tpr][tp]
      for (let t = 0; t < tpv.length; t++) { // loop over members of target partition
        const tv = tpv[t]
        const tx = tv.x + tp
        const ty = tv.y + tpr
        let fx = 0
        let fy = 0
        for (let spr = tpr - 1; spr < tpr + 2; spr++) { // loop over source partitions rows
          for (let sp = tp - 1; sp < tp + 2; sp++) { // loop over source partitions
            const spv = arena[(spr + arenaHeight) % arenaHeight][(sp + arenaWidth) % arenaWidth] // TRICK!
            for (let s = 0; s < spv.length; s++) { // loop over members of source partition
              if (tpr === spr && tp === sp && t === s) {
                continue // do not interact with itself
              }
              const sv = spv[s]
              const sx = sv.x + sp // TRICK!
              const sy = sv.y + spr
              const dx = sx - tx
              const dy = sy - ty
              const r = Math.hypot(dx, dy)
              var f = 0
              if (r < r1) {
                f = (-1 + r / r1) * 100
              } else if (r < r2) {
                f = forceMatrix[tv.tp][sv.tp] * (r - r1) / (r2 - r1) // TODO: pre calculate denominator
              } else if (r < 1) {
                f = forceMatrix[tv.tp][sv.tp] * (1 - r) / (1 - r2)
              }
              fx += f * dx / r
              fy += f * dy / r
            }
          }
        }
        fx -= tv.vx * .8 // TODO
        fy -= tv.vy * .8
        tv.vx += dt * fx
        tv.vy += dt * fy
      }
    }
  }
}

function moveArena(dt) {
  const a = createEmptyArena()
  for (let pr = 0; pr < arenaHeight; pr++) { // loop over partitions rows
    for (let p = 0; p < arenaWidth; p++) { // loop over partitions
      const pv = arena[pr][p]
      for (let i = 0; i < pv.length; i++) { // loop over members of partition
        const v = pv[i]
        const x = (p + v.x + dt * v.vx + arenaWidth) % arenaWidth
        const y = (pr + v.y + dt * v.vy + arenaHeight) % arenaHeight
        v.x = x % 1
        v.y = y % 1
        a[Math.floor(y)][Math.floor(x)].push(v)
      }
    }
  }
  arena = a
}

let prevMs = 0
function animate(ms) {
  if (prevMs === 0) {
    prevMs = ms
  }
  let dt = (ms - prevMs) / 1000 * 20 // TODO time factor
  if (dt !== 0) { // TODO if dt is too big, skip too?
    if (dt > .015) {
      dt = .015
    }
    accelerate(dt)
    moveArena(dt)
    drawArena()
    prevMs = ms
  }
  window.requestAnimationFrame(animate)
}

window.requestAnimationFrame(animate)

/*
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
*/
