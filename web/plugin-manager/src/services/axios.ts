import React from 'react'
import axios, { AxiosError, AxiosResponse } from 'axios'
import { Plugin } from '../interfaces/IPlugin'

export class axiosServices {
  constructor() {
  
  }

  static manageLifePlugins (ID:number, command:string) : Promise<AxiosResponse<any>> {
    return axios.post<any>(`${window.location.hostname}/plugin-manager/${ID}/execute`, {
      command: command,
    })
  }

  static getPluginsInfo () : Promise<AxiosResponse<Plugin[]>> {
    return  axios.get<Plugin[]>(`${window.location.hostname}/plugin-manager`);
  }

  static getpluginInfo (ID:number) : Promise<AxiosResponse<Plugin, any>> {
    return  axios.get<Plugin>(`${window.location.hostname}/plugin-manager/${ID}`);      
  }

  static  deletePlugins (ID:number) : Promise<AxiosResponse<any>> {
    return axios.delete<any>(`${window.location.hostname}/plugin-manager/${ID}`)

    }

  static uploadPlugins (fileName:string) : string {
    // const formData = new FormData()
    // formData.append('file', file)
    // return axios.post<any>(`${window.location.hostname}/plugin-manager/upload`, formData)
    return "Not implemented yet"
    }
  }

export default axiosServices