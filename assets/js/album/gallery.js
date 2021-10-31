$(function () {
    window.lightGallery(
      document.getElementById("album-gallery"),
      {
        autoplayFirstVideo: false,
        pager: true,
        galleryId: "photos",
        plugins: [lgThumbnail],
        licenseKey: '0000-0000-000-0000',
        mobileSettings: {
          controls: true,
          showCloseIcon: true,
          download: true,
          rotate: true
        }
      }
    );
});

