import plusSmall from '../assets/pictos/plusSmall.svg';
import PlusSmallWhite from '../assets/pictos/PlusSmallWhite.svg';
import { UIStore } from '../store/UIStore';
const ManagePluginCard = () => {
  const plus = UIStore.useState((s) =>
    s.theme == 'light' ? plusSmall : PlusSmallWhite,
  );
  return (
    <div
      onClick={() => {
        window.open('/thyra/plugin-manager');
      }}
      className="flex flex-col justify-center items-start p-10 gap-5 w-64 h-72 rounded-2xl"
    >
      <button>
        <img
          src={plus}
          alt="plusSmall"
          className="rounded-full w-16 h-16 mx-auto bg-bgCard hover:bg-hoverbgCard"
        />
        <h2 className="label text-font">Manage plugin</h2>
      </button>
    </div>
  );
};

export default ManagePluginCard;
