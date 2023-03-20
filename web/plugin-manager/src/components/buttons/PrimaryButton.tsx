import React from 'react'
import { UIStore } from '../../store/UIStore';

type Props = {
  label : string,
  onClick : () => void,
  isprimary? : boolean,
  iconPathLight? : string,
  iconPathDark? : string,
}

const PrimaryButton = (props: Props) => {
  const isThemeLight = UIStore.useState(s => (s.theme == "light" ? true : false));
  return (

      <button>
        <div className="flex flex-row justify-center items-center gap-2 w-28 h-12 bg-primaryButton border-[1px] border-solid border-border rounded-md cursor-pointer hover:bg-hoverSecondaryButton ">
          <p className="text-invertedfont">{props.label}</p>
          {props.iconPathLight ? <img className='w-2 h-2 mt-1' src={isThemeLight ? props.iconPathLight : props.iconPathDark } alt="No Icon Found" /> : <></> }
        </div>
      </button>

  )
}

export default PrimaryButton