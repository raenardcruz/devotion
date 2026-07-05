export interface MysteryVisual {
  image: string;
  meditation: string;
}

export const MYSTERY_VISUALS: Record<string, MysteryVisual> = {
  // Joyful
  "The Annunciation": {
    image: "/images/rosary/annunciation.png", // Classical Annunciation painting style
    meditation: "Contemplate the humility of Mary and her total surrender to the divine will. Let the silence of this moment permeate your prayer."
  },
  "The Visitation": {
    image: "/images/rosary/visitation.png",
    meditation: "Reflect on the joy of charity and the promptness of Mary to serve her cousin Elizabeth. Seek to bring Christ to those you meet today."
  },
  "The Nativity": {
    image: "/images/rosary/annunciation.png",
    meditation: "Ponder the poverty and humility of the stable in Bethlehem. Let the light of the newborn Savior shine in the quiet places of your heart."
  },
  "The Presentation": {
    image: "/images/rosary/presentation.png",
    meditation: "Offer yourself entirely to God in imitation of Mary and Joseph presenting the Infant Jesus. Pray for the spirit of obedience and purity."
  },
  "The Finding in the Temple": {
    image: "/images/rosary/finding_temple.png",
    meditation: "Seek Jesus with the same diligence when He seems distant. Rest in the knowledge that He is always about His Father's business."
  },
  // Luminous
  "The Baptism of Jesus": {
    image: "/images/rosary/annunciation.png",
    meditation: "Hear the voice of the Father declaring you as His beloved child. Renew your baptismal call to live in holiness and justice."
  },
  "The Wedding at Cana": {
    image: "/images/rosary/wedding_cana.png",
    meditation: "Do whatever He tells you. Trust in the quiet intercession of Mary, who notices every human need and points us to her Son."
  },
  "The Proclamation of the Kingdom": {
    image: "/images/rosary/transfiguration.png",
    meditation: "Repent and believe in the Gospel. Open your heart to the transformative power of God's mercy and extend it to those around you."
  },
  "The Transfiguration": {
    image: "/images/rosary/annunciation.png",
    meditation: "Gaze upon the radiant glory of Christ. Let the hope of future glory strengthen you to carry your daily crosses with faith."
  },
  "The Institution of the Eucharist": {
    image: "/images/rosary/institution_eucharist.png",
    meditation: "Receive the Gift of Love. Rest in the real presence of Jesus, who feeds us with His own Body and Blood to remain in us."
  },
  // Sorrowful
  "The Agony in the Garden": {
    image: "/images/rosary/agony_garden.png",
    meditation: "Not my will, but Yours be done. Stand with Jesus in His darkest hour of sorrow, offering Him your console and obedience."
  },
  "The Scourging at the Pillar": {
    image: "/images/rosary/scourging_pillar.png",
    meditation: "Offer your physical trials and suffering in union with the scourged body of Christ, who bore our wounds for our salvation."
  },
  "The Crowning with Thorns": {
    image: "/images/rosary/crowning_thorns.png",
    meditation: "Ponder the mockery faced by the King of Kings. Pray for humility to accept insults and the grace to crown Him Lord of your thoughts."
  },
  "The Carrying of the Cross": {
    image: "/images/rosary/carrying_cross.png",
    meditation: "Take up your cross daily and follow Him. Like Simon of Cyrene, be ready to help lift the burdens of others along the way."
  },
  "The Crucifixion": {
    image: "/images/rosary/annunciation.png",
    meditation: "Father, forgive them, for they know not what they do. Stand at the foot of the Cross with Mary and receive the depth of His saving love."
  },
  // Glorious
  "The Resurrection": {
    image: "/images/rosary/resurrection.png",
    meditation: "He is risen! Let the victory of Christ over death fill you with hope. Live as a witness to the resurrection in your daily actions."
  },
  "The Ascension": {
    image: "/images/rosary/ascension.png",
    meditation: "Set your eyes on the things of heaven. Trust that Jesus has gone to prepare a place for us and remains with us until the end of time."
  },
  "The Descent of the Holy Spirit": {
    image: "/images/rosary/descent_holy_spirit.png",
    meditation: "Come, Holy Spirit, fill the hearts of Your faithful. Pray for a renewal of the spiritual gifts and the fire of charity in your life."
  },
  "The Assumption of Mary": {
    image: "/images/rosary/assumption.png",
    meditation: "Look to Mary, taken body and soul into heaven, as a sign of sure hope and solace. Pray for a happy death in the grace of God."
  },
  "The Coronation of Mary": {
    image: "/images/rosary/coronation.png",
    meditation: "Rejoice in the crowning of our Queen and Mother. Cast all your cares upon her maternal heart, knowing she is your advocate."
  }
};

export const getMysteryVisual = (title: string): MysteryVisual => {
  return MYSTERY_VISUALS[title] || {
    image: "https://images.unsplash.com/photo-1580136579312-94651dfd596d?q=80&w=600&auto=format&fit=crop",
    meditation: "Let the silence of this moment permeate your prayer."
  };
};
