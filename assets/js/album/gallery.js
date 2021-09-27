$(function () {
    window.lightGallery(
      document.getElementById("album-gallery"),
      {
        autoplayFirstVideo: false,
        pager: false,
        galleryId: "nature",
        plugins: [lgThumbnail],
        licenseKey: '0000-0000-000-0000',
        mobileSettings: {
          controls: false,
          showCloseIcon: false,
          download: false,
          rotate: false
        }
      }
    );
});

