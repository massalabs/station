import axios, { AxiosResponse } from 'axios'
import { Plugin } from '../../../shared/interfaces/IPlugin'

export class axiosServices {

  //Developement process to get the plugin-manager API in local environment
  //1. Run Thyra
  //4. Run the plugin-manager UI in local environment

  static LocalPort = () => {
    return "http://localhost:4200"
  }

  

  static localinstallPlugin(url: string): Promise<AxiosResponse<number>> {
    return axios.post<any>(`${this.LocalPort}/plugin-manager?source=${url}`)
  }


//

  static installPlugin(url: string): Promise<AxiosResponse<number>> {
    return axios.post<any>(`/plugin-manager?source=${url}`)
  }

  static manageLifePlugins(ID: string, command: string): Promise<AxiosResponse<number>> {
    return axios.post<any>(`/plugin-manager/${ID}/execute`, {
      command: command,
    })
  }

  static getPluginsInfo(): Promise<AxiosResponse<Plugin[]>> {
    return axios.get<Plugin[]>(`/plugin-manager`);
  }

  static getpluginInfo(ID: string): Promise<AxiosResponse<{ status: string }>> {
    return axios.get<{ status: string }>(`/plugin-manager/${ID}`);
  }

  static deletePlugins(ID: string): Promise<AxiosResponse<any>> {
    return axios.delete<any>(`/plugin-manager/${ID}`)

  }

  static uploadPlugins(fileName: string): string {
    // const formData = new FormData()
    // formData.append('file', file)
    // return axios.post<any>(`${window.location.hostname}/plugin-manager/upload`, formData)
    return "Not implemented yet"
  }
}

export default axiosServices