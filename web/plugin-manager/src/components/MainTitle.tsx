const MainTitle = (title: { title: string }) => {
    return (
        <p className=" display flex-row flex text-font">
            <p className="text-brand">↳</p> {title.title}
        </p>
    );
};

export default MainTitle;
