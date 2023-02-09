import React from 'react'
import axios, { AxiosError, AxiosResponse } from 'axios'
import { Plugin } from '../../../shared/interfaces/IPlugin'

export class axiosServices {
  constructor() {
  
  }

  static manageLifePlugins (ID:string, command:string) : Promise<AxiosResponse<number>> {
    return axios.post<any>(`/plugin-manager/${ID}/execute`, {
      command: command,
    })
  }

  static getPluginsInfo () : Promise<AxiosResponse<Plugin[]>> {
    return  axios.get<Plugin[]>(`/plugin-manager`);
  }

  static getpluginInfo (ID:string) : Promise<AxiosResponse<string, any>> {
    console.log("getpluginInfo" + ID)
    return  axios.get<string>(`/plugin-manager/${ID}`);      
  }

  static  deletePlugins (ID:string) : Promise<AxiosResponse<any>> {
    return axios.delete<any>(`/plugin-manager/${ID}`)

    }

  static uploadPlugins (fileName:string) : string {
    // const formData = new FormData()
    // formData.append('file', file)
    // return axios.post<any>(`${window.location.hostname}/plugin-manager/upload`, formData)
    return "Not implemented yet"
    }
  }

export default axiosServices