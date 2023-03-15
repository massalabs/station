import React from 'react';
import Arrow6 from '../assets/pictos/arrow6.svg';
type Props = {
  logo: string;
  name: string;
  description: string;
};

export const PluginCard = (props: Props) => {
  return (
      <div className="flex flex-col justify-center items-start p-10 gap-5 w-64 h-72 border-[1px] border-solid border-black rounded-2xl bg-bgCard">
        <img src={props.logo} alt="Album" className="rounded-3xl w-16 h-16" />
        <h2 className='label2'>{props.name}</h2>
        <p className='text2 overflow-hidden whitespace-pre-wrap max-w-full'>{props.description}</p>
        <button>
          <img src={Arrow6} alt="Album" className="rounded-3xl w-6 h-6" />
        </button>
      </div>
  );
};
