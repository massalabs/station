import React, { useState } from 'react'
import { axiosServices } from '../services/axios';
import PrimaryButton from './buttons/PrimaryButton';
import SecondaryButton from './buttons/SecondaryButton';

type Props  = {
    callbackToParent: (data: string) => void;
    label:string;
    placeholder: string;
    buttonValue: string;
    axiosCall: (data: string) => void;
    error?:string;
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
        <p className='text-font my-1'>{props.label}</p>
        <input type="text" className="text-font w-full mb-4 rounded-md bg-primaryBG" placeholder={props.placeholder} onChange={handleInputValueChange} />
        <PrimaryButton label={props.buttonValue} onClick={handleSubmit} width={" w-full"}/>
    </div>
  )
}

export default LabelButton