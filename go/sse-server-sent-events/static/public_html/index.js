'use strict'

const eForm = document.getElementsByTagName('form')[0]
const eInput = document.getElementsByTagName('input')[0]
const eAllert = document.getElementsByTagName('div')[0]

function updateAlert() {
  const text = eInput.value
  const room = text.replaceAll(/[^0-9a-zA-Z_-]+/g, '')
  if (room !== text) {
    eAllert.style.backgroundColor = '#ff0'
  } else {
    eAllert.style.backgroundColor = ''
  }
}

eInput.onkeyup = updateAlert
eForm.onsubmit = (e) => {
  e.preventDefault()
  const text = eInput.value
  const room = text.replaceAll(/[^0-9a-zA-Z_-]+/g, '')
  if (room === text && room !== '') {
    location.href = '/' + room
  } else {
    eInput.value = room
    updateAlert()
  }
}
