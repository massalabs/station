import React from 'react'
import { UIStore } from '../../store/UIStore';

type Props = {
  label : string,
  onClick : () => void,
  width? : string,
}

const SecondaryButton = (props: Props) => {
  return (
      <button>
        <div className={"flex flex-row justify-center items-center gap-2 h-12 bg-secondaryButton border-[1px] border-solid border-border rounded-md cursor-pointer " +( props.width ?? "w-28" ) + " hover:bg-hoverPrimaryButton" }>
          <p className={"text-font"}>{props.label}</p>
        </div>
      </button>
  )
}

export default SecondaryButton