import { useState } from 'react'
import loop from './assets/loop.wav'
import './App.css'
import * as Tone from "tone"

function App() {
  const [isPlaying, setIsPlaying] = useState(false)

  const handleClick = () => {
    if (isPlaying) {
      return
    }
    const fetchData = async () => {
      const response = await fetch('https://decay.oren.cool')
      const data = await response.json()
      return data
    }

    fetchData().then((res) => {
      setIsPlaying(true)
      // Setup Tone.JS
      Tone.start()

      const player = new Tone.Player({
        url: loop,
        autostart: true,
        loop: true,
      }) 


      // Note to self: I don't know why I had the decay value inverted, but I was too lazy to change it.
      const mappedValue = 1 - res.value

      const distortion = new Tone.Distortion(mappedValue).toDestination()
      const reverb = new Tone.Reverb(mappedValue * 1000000).toDestination()
      const delay = new Tone.FeedbackDelay(mappedValue, mappedValue).toDestination()

      player.connect(delay).connect(reverb).connect(distortion)
    })
  }

  return (
    <div className="container">
      <div className='text-container'>
      <h1>
        This sound is impermanent.
      </h1>

      <h1>
        It will decay over time.
      </h1>
      <p> In one year, it will be gone.</p>

      </div>


      <button onClick={() => handleClick()}>
        Listen
      </button>
    </div>
  )
}

export default App
