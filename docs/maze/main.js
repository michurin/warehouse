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

function mod(n, m) {
  return ((n % m) + m) % m;
}

function randomDirs() {
  const x = [{ x: 1, y: 0 }, { x: -1, y: 0 }, { x: 0, y: 1 }, { x: 0, y: -1 }]
  var r, n = x.length
  while (n > 0) {
    r = Math.floor(Math.random() * n)
    n--
    [x[n], x[r]] = [x[r], x[n]]
  }
  return x
}

function mazePut(m, gen, xo, yo) {
  m[xo][yo] = gen
  const g = ++gen // gen+1 for gate-cell
  gen++ // gen+2 for next
  const dd = randomDirs()
  dd.forEach(d => {
    var x = mod(xo + d.x * 2, m.length)
    var y = mod(yo + d.y * 2, m[0].length)
    if (m[x][y] === undefined) {
      var x1 = mod(xo + d.x, m.length)
      var y1 = mod(yo + d.y, m[0].length)
      m[x1][y1] = g
      mazePut(m, gen, x, y)
    }
  })
}

function arr2d(size) {
  const m = new Array(size)
  for (var i = 0; i < m.length; i++) { m[i] = new Array(size) }
  return m
}

function mazeGen() {
  const pmaze = arr2d(16)
  mazePut(pmaze, 0, 0, 0)
  const maze = arr2d(16 * 4)
  for (var j = 0; j < maze.length; j++) {
    for (var i = 0; i < maze[j].length; i++) {
      const q = Math.floor(j / 4)
      const p = Math.floor(i / 4)
      const q1 = Math.floor((j + 2) / 4) % 16
      const p1 = Math.floor((i + 2) / 4) % 16
      maze[j][i] = (
        pmaze[q][p] === undefined &&
        pmaze[q][p1] === undefined &&
        pmaze[q1][p] === undefined &&
        pmaze[q1][p1] === undefined
      ) ? 100001 : undefined
    }
  }
  mazePut(maze, 0, 0, 0)
  for (var j = 0; j < maze.length; j++) {
    for (var i = 0; i < maze[j].length; i++) {
      if (maze[j][i] > 100000) {
        maze[j][i] = undefined
      }
    }
  }
  mx = 0
  maze.forEach(s => {
    s.forEach(e => {
      if (e === undefined) {
        return
      }
      if (e > mx) {
        mx = e
      }
    })
  })
  console.log(mx)
  mx++
  for (var j = 0; j < maze.length; j++) {
    var p = maze[j]
    for (var i = 0; i < p.length; i++) {
      if (p[i] === undefined) {
        continue
      }
      p[i] /= mx
    }
  }
  return maze
}

function color(m, s) {
  if (m === undefined) {
    return '#000'
  }
  return `rgb(0, ${Math.floor(50 + m * 200)}, 0)`
  return '#fff'
  //m *= 5
  m += s
  m -= Math.floor(m)
  return m > .9 ? '#fff' : '#000'
  m = 1 - m
  if (m < 0) {
    m = 0
  }
  if (m > 1) {
    m = 1
  }
  return `rgb(0, ${Math.floor(50 + m * 200)}, 0)`
}

function drawMoze(maze) {
  const base = 400

  const canvas = document.getElementById('a')
  canvas.style.width = '400px'
  canvas.style.height = '400px'
  canvas.width = Math.round(base)
  canvas.height = Math.round(base)
  const context = canvas.getContext('2d')
  context.fillStyle = '#000'
  context.fillRect(0, 0, canvas.width, canvas.height)

  for (var j = 0; j < maze.length; j++) {
    for (var i = 0; i < maze[j].length; i++) {
      context.fillStyle = color(maze[j][i], 0)
      context.fillRect(j * 4, i * 4, 3, 3)
    }
  }
}

function drawInvMoze(maze, colorShift) {
  const base = 800

  const canvas = document.getElementById('b')
  canvas.style.width = '400px'
  canvas.style.height = '400px'
  canvas.width = Math.round(base)
  canvas.height = Math.round(base)
  const context = canvas.getContext('2d')
  context.fillStyle = '#000'
  context.fillRect(0, 0, canvas.width, canvas.height)

  for (var j = 0; j < base; j++) {
    for (var i = 0; i < base; i++) {
      const x = i - base / 2
      const y = j - base / 2
      const a = Math.atan2(y, x) / 2 / Math.PI * maze.length
      const rr = (x * x + y * y)
      if (rr === 0) {
        continue
      }
      const x1 = Math.log(Math.sqrt(rr)) * 40
      const y1 = a * 3
      const m = maze[mod(Math.floor(y1), maze.length)][mod(Math.floor(x1), maze[0].length)]
      context.fillStyle = color(m, colorShift)
      context.fillRect(j, i, 1, 1)
    }
  }
}

const maze = mazeGen()
console.log(maze)
drawMoze(maze)

function D(s) {
  drawInvMoze(maze, s)
  console.log(s)
  //  s += .1
  //  setTimeout(function() { D(s) }, 10)
}
D(0)
