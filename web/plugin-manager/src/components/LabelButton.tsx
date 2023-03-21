import React, { useState } from 'react'
import { axiosServices } from '../services/axios';
import SecondaryButton from './buttons/SecondaryButton';

type Props  = {
    callbackToParent: (data: string) => void;
    label:string;
    placeholder: string;
    buttonValue: string;
    axiosCall: (data: string) => void;
}

const LabelButton = (props: Props) => {
    const [value, setValue] = useState('');

    function handleInputValueChange(event: any) {
        setValue(event.target.value);
    }

    async function handleSubmit() {
        props.axiosCall(value)
        props.callbackToParent
    }

  return (
    <div className=''>
        <p className='text-font my-2'>{props.label}</p>
        <input type="text" className="w-full mb-2" placeholder={props.placeholder} onChange={handleInputValueChange} />
        <SecondaryButton label={props.buttonValue} onClick={handleSubmit} width={"w-full"}/>
    </div>
  )
}

export default LabelButton