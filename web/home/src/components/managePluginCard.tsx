import React from 'react'
import plusSmall from '../../assets/pictos/plusSmall.svg'

const managePluginCard = () => {
  return (
    <div className="flex flex-col justify-center items-start p-10 gap-5 w-64 h-72 rounded-2xl">
    <button>
    <img src={plusSmall} alt="plusSmall" className="rounded-3xl w-16 h-16" />
    <h2 className='label'>Manage plugin</h2>
    </button>
  </div>
  )
}

export default managePluginCard