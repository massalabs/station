import React from 'react'

type Props = {
    name:string,
    logo:string,
    description: string,
    version: string,
    online: boolean,
}

function pluginBlock({}: Props) {
  return (
    <div>pluginBlock</div>
  )
}

export default pluginBlock