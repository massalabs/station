import React from 'react'
import { UIStore } from '../../store/UIStore';

type Props = {
  label : string,
  onClick : () => void,
  isprimary? : boolean,
  iconPathLight? : string,
  iconPathDark? : string,
  width? : string,
}

const PrimaryButton = (props: Props) => {
  return (
      <button>
        <div className={"flex flex-row justify-center items-center gap-2 h-12 bg-secondaryButton border-[1px] border-solid border-border rounded-md cursor-pointer hover:bg-hoverPrimaryButton " + props.width ?? " w-28"}>
          <p className={props.isprimary ? "text-invertedfont" : "text-font"}>{props.label}</p>
        </div>
      </button>
  )
}

export default PrimaryButton