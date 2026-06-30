export interface MysteryVisual {
  image: string;
  meditation: string;
}

export const MYSTERY_VISUALS: Record<string, MysteryVisual> = {
  // Joyful
  "The Annunciation": {
    image: "https://images.unsplash.com/photo-1580136579312-94651dfd596d?q=80&w=600&auto=format&fit=crop", // Classical Annunciation painting style
    meditation: "Contemplate the humility of Mary and her total surrender to the divine will. Let the silence of this moment permeate your prayer."
  },
  "The Visitation": {
    image: "https://images.unsplash.com/photo-1601887389937-0b02c26b6c3c?q=80&w=600&auto=format&fit=crop",
    meditation: "Reflect on the joy of charity and the promptness of Mary to serve her cousin Elizabeth. Seek to bring Christ to those you meet today."
  },
  "The Nativity": {
    image: "https://images.unsplash.com/photo-1543872084-c7bd3822856f?q=80&w=600&auto=format&fit=crop",
    meditation: "Ponder the poverty and humility of the stable in Bethlehem. Let the light of the newborn Savior shine in the quiet places of your heart."
  },
  "The Presentation": {
    image: "https://images.unsplash.com/photo-1549429712-40f4244fa7be?q=80&w=600&auto=format&fit=crop",
    meditation: "Offer yourself entirely to God in imitation of Mary and Joseph presenting the Infant Jesus. Pray for the spirit of obedience and purity."
  },
  "The Finding in the Temple": {
    image: "https://images.unsplash.com/photo-1548625361-155deee223d0?q=80&w=600&auto=format&fit=crop",
    meditation: "Seek Jesus with the same diligence when He seems distant. Rest in the knowledge that He is always about His Father's business."
  },
  // Luminous
  "The Baptism of Jesus": {
    image: "https://images.unsplash.com/photo-1515263487990-61b07816b324?q=80&w=600&auto=format&fit=crop",
    meditation: "Hear the voice of the Father declaring you as His beloved child. Renew your baptismal call to live in holiness and justice."
  },
  "The Wedding at Cana": {
    image: "https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?q=80&w=600&auto=format&fit=crop",
    meditation: "Do whatever He tells you. Trust in the quiet intercession of Mary, who notices every human need and points us to her Son."
  },
  "The Proclamation of the Kingdom": {
    image: "https://images.unsplash.com/photo-1447069387593-a5de0862481e?q=80&w=600&auto=format&fit=crop",
    meditation: "Repent and believe in the Gospel. Open your heart to the transformative power of God's mercy and extend it to those around you."
  },
  "The Transfiguration": {
    image: "https://images.unsplash.com/photo-1451187580459-43490279c0fa?q=80&w=600&auto=format&fit=crop",
    meditation: "Gaze upon the radiant glory of Christ. Let the hope of future glory strengthen you to carry your daily crosses with faith."
  },
  "The Institution of the Eucharist": {
    image: "https://images.unsplash.com/photo-1544427920-c49ccfb85579?q=80&w=600&auto=format&fit=crop",
    meditation: "Receive the Gift of Love. Rest in the real presence of Jesus, who feeds us with His own Body and Blood to remain in us."
  },
  // Sorrowful
  "The Agony in the Garden": {
    image: "https://images.unsplash.com/photo-1509021436665-8f37da187925?q=80&w=600&auto=format&fit=crop",
    meditation: "Not my will, but Yours be done. Stand with Jesus in His darkest hour of sorrow, offering Him your console and obedience."
  },
  "The Scourging at the Pillar": {
    image: "https://images.unsplash.com/photo-1594847067869-0855c843f412?q=80&w=600&auto=format&fit=crop",
    meditation: "Offer your physical trials and suffering in union with the scourged body of Christ, who bore our wounds for our salvation."
  },
  "The Crowning with Thorns": {
    image: "https://images.unsplash.com/photo-1544816155-12df9643f363?q=80&w=600&auto=format&fit=crop",
    meditation: "Ponder the mockery faced by the King of Kings. Pray for humility to accept insults and the grace to crown Him Lord of your thoughts."
  },
  "The Carrying of the Cross": {
    image: "https://images.unsplash.com/photo-1553531384-397c80973a0b?q=80&w=600&auto=format&fit=crop",
    meditation: "Take up your cross daily and follow Him. Like Simon of Cyrene, be ready to help lift the burdens of others along the way."
  },
  "The Crucifixion": {
    image: "https://images.unsplash.com/photo-1507054798601-57f1fa70446b?q=80&w=600&auto=format&fit=crop",
    meditation: "Father, forgive them, for they know not what they do. Stand at the foot of the Cross with Mary and receive the depth of His saving love."
  },
  // Glorious
  "The Resurrection": {
    image: "https://images.unsplash.com/photo-1518005020951-eccb494ad742?q=80&w=600&auto=format&fit=crop",
    meditation: "He is risen! Let the victory of Christ over death fill you with hope. Live as a witness to the resurrection in your daily actions."
  },
  "The Ascension": {
    image: "https://images.unsplash.com/photo-1469474968028-56623f02e42e?q=80&w=600&auto=format&fit=crop",
    meditation: "Set your eyes on the things of heaven. Trust that Jesus has gone to prepare a place for us and remains with us until the end of time."
  },
  "The Descent of the Holy Spirit": {
    image: "https://images.unsplash.com/photo-1447752875215-b2761acb3c5d?q=80&w=600&auto=format&fit=crop",
    meditation: "Come, Holy Spirit, fill the hearts of Your faithful. Pray for a renewal of the spiritual gifts and the fire of charity in your life."
  },
  "The Assumption of Mary": {
    image: "https://images.unsplash.com/photo-1501854140801-50d01698950b?q=80&w=600&auto=format&fit=crop",
    meditation: "Look to Mary, taken body and soul into heaven, as a sign of sure hope and solace. Pray for a happy death in the grace of God."
  },
  "The Coronation of Mary": {
    image: "https://images.unsplash.com/photo-1563245372-f21724e3856d?q=80&w=600&auto=format&fit=crop",
    meditation: "Rejoice in the crowning of our Queen and Mother. Cast all your cares upon her maternal heart, knowing she is your advocate."
  }
};

export const getMysteryVisual = (title: string): MysteryVisual => {
  return MYSTERY_VISUALS[title] || {
    image: "https://images.unsplash.com/photo-1580136579312-94651dfd596d?q=80&w=600&auto=format&fit=crop",
    meditation: "Let the silence of this moment permeate your prayer."
  };
};
