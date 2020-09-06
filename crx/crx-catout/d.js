(() => {
  setInterval(() => {
    [
      'ad_topslot_a',
      'ad_topslot_b',
      'ad_leftslot',
      'ad_rightslot',
      'ad_rightslot2',
      'ad_houseslot_a',
      'ad_houseslot_b',
      'ad_contentslot_1',
      'ad_contentslot_2',
      'ad_contentslot_3',
      'ad_contentslot_4',
      'ad_btmslot_a',
    ].forEach(x => {
      const e = document.getElementById(x);
      if (e) {
        e.remove();
      }
    });
    document.querySelectorAll('.loginPopupContainer').forEach(x => {
      x.remove();
    });
  }, 2000);
})();
