import React, { useState } from 'react'
import halfArrow from '../assets/pictos/halfArrow.svg'
import LabelButton from './LabelButton'
type Props = {
    label2:  string
    text3 : string
    propsLabelButton: {
        callbackToParent: (data: string) => void;
        label:string;
        placeholder: string;
        buttonValue: string;
        axiosCall: (data: string) => void;
    }
}

const SmallCardExtended = (props: Props) => {
    const [isExtended, setIsExtended] = useState(false)
    const HandleOnClickExtend = () => {
        setIsExtended(!isExtended)
        //Change the className of the container to grow the height of the section
        return " h-56 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard"
    }
    const sectionClassName = () => {
        if(isExtended){
            return "h-56 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard hover:bg-hoverbgCard"
        }
        else{
            return "h-28 w-80 max-w-lg p-6 gap-3 border-[1px] border-solid border-border rounded-2xl bg-bgCard hover:bg-hoverbgCard"
        }
    } 
  return (
    <section className={sectionClassName()}>
        {/* Definition de la card  */}
            <div className='flex flex-row '>
            <div className='flex flex-col gap-3'>
                <p className='label2 text-font'>
                    Install a plugin
                </p>
                <p className='text3 text-font'>
                    Install a plugin using .zip URL
                </p>
            </div>
            <div className='flex self-center mx-auto'>
                <img className='w-8 h-4 hover:animate-rotate180Smooth' src={halfArrow} alt="" onClick={HandleOnClickExtend} />
                {/* contains the icon to grow the container */}
            </div>
            </div>
            {isExtended && <div className='flex flex-col mt-3'>
                <LabelButton {...props.propsLabelButton}/>
                </div>}
    </section>
  )
}

export default SmallCardExtended