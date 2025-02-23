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

// For developers: how to check
// eslint --no-eslintrc --fix main.js

// SETTINGS [1/3]: Geometry

const partitionSize = 2
const partitioningFactor = 10

const arenaWidth = 400 // in partitions
const arenaHeight = 400

// SETTINGS [2/3]: Dynamics of particles

const initialVelocity = 1 // partitionSize * partitioningFactor / second

const numgerOfParticals = 400 // for each individual specie

const r1 = .4 // radius of repulsion; must be less than one
const r2 = .7 // radius of interaction according matrix; must be less than one and greater than r1
const forceMatrix = [ // N-by-N square matrix, where N is number of species; positive — attraction, negative — repulsion
  [3, .5],
  [-.5, 3],

  //  [.4, .2, 0, 0],
  //  [-.05, .4, 0, 0],
  //  [0, 0, .4, .2],
  //  [0, 0, -.005, .4],
]

const dynamicViscosity = 2
const repulsionScale = 20

// SETTINGS [3/3]: Dynamics

const integrationTimeFactor = 4 // bigger — faster and less accuracy, smaller — slower and better accuracy
const integrationIntervalLimit = .02 // hard limit of integration interval in seconds

// END OF SETTINGS

const arenaWidthModAlign = arenaWidth * 100 // is stupid helper for true m = ((x%n)+n)%n approach
const arenaHeightModAlign = arenaHeight * 100

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

function color(n, total) { // function heights where colors come from
  return `hsl(${Math.round(n * 360 / total)} 100% 50%)`
}

function putPartical(x, y, v, va, tp, clr) { // function heights how partitioning works
  arena[Math.floor(y)][Math.floor(x)].push({
    x: x % 1, // coordinates, related partition
    y: y % 1,
    vx: v * Math.cos(va),
    vy: v * Math.sin(va),
    tp: tp,
    clr: clr,
  })
}

for (let tp = 0; tp < forceMatrix.length; tp++) {
  const clr = color(tp, forceMatrix.length)
  for (let s = 0; s < numgerOfParticals; s++) {
    putPartical(arenaWidth * Math.random(), arenaHeight * Math.random(), initialVelocity, Math.random() * Math.PI * 2, tp, clr)
  }
}

// ---- setup canvas

const integrationIntervalElement = document.getElementById('integration-interval')

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
        /*
                context.beginPath()
                context.arc((p + v.x) * partitionSize, (pr + v.y) * partitionSize, r1 * partitionSize * partitioningFactor, 0, 2 * Math.PI);
                context.fillStyle = 'rgba(255, 255, 255, .03)'
                context.fill();
                context.closePath()

                context.beginPath()
                context.arc((p + v.x) * partitionSize, (pr + v.y) * partitionSize, partitionSize * partitioningFactor, 0, 2 * Math.PI);
                context.fillStyle = 'rgba(255, 255, 0, .03)'
                context.fill();
                context.closePath()
        */
        context.fillStyle = v.clr
        context.fillRect((p + v.x) * partitionSize - 1, (pr + v.y) * partitionSize - 1, 3, 3)
      }
    }
  }
}

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
        for (let spr = tpr - partitioningFactor; spr <= tpr + partitioningFactor; spr++) { // loop over source partitions rows
          for (let sp = tp - partitioningFactor; sp <= tp + partitioningFactor; sp++) { // loop over source partitions
            const spv = arena[(spr + arenaHeightModAlign) % arenaHeight][(sp + arenaWidthModAlign) % arenaWidth] // TRICK! Part one: how we takes partition
            for (let s = 0; s < spv.length; s++) { // loop over members of source partition
              if (tpr === spr && tp === sp && t === s) {
                continue // do not interact with itself
              }
              const sv = spv[s]
              const sx = sv.x + sp // TRICK! Part two: how we calculate distances (without %). It does wrapping arena on bounds
              const sy = sv.y + spr
              const dx = sx - tx
              const dy = sy - ty
              const r = Math.hypot(dx, dy) / partitioningFactor
              let f = 0
              if (r < r1) {
                f = (-1 + r / r1) * repulsionScale
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
        fx -= tv.vx * dynamicViscosity
        fy -= tv.vy * dynamicViscosity
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
        const x = (p + v.x + dt * v.vx + arenaWidthModAlign) % arenaWidth
        const y = (pr + v.y + dt * v.vy + arenaHeightModAlign) % arenaHeight
        v.x = x % 1
        v.y = y % 1
        a[Math.floor(y)][Math.floor(x)].push(v)
      }
    }
  }
  arena = a
}

let animationCount = 0
let cummulativeAnimationInterval = 0
let prevMs = 0
function animate(ms) {
  if (prevMs === 0) {
    prevMs = ms
  }
  let dt = ms - prevMs
  if (dt !== 0) {
    animationCount++
    dt /= 1000 // seconds
    if (dt > integrationIntervalLimit) {
      dt = integrationIntervalLimit
    }
    cummulativeAnimationInterval = cummulativeAnimationInterval * .99 + dt
    integrationIntervalElement.innerText = (['/', '-', '\\', '|'][animationCount % 4]) + ' ' + dt.toFixed(5) + ' ' + (cummulativeAnimationInterval / 100).toFixed(5)
    dt *= integrationTimeFactor
    accelerate(dt)
    moveArena(dt)
    drawArena()
    prevMs = ms
  }
  window.requestAnimationFrame(animate)
}

window.requestAnimationFrame(animate)
